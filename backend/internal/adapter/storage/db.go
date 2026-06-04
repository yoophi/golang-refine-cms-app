package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"      // postgres 드라이버 등록 (이름: "postgres")
	_ "modernc.org/sqlite"     // 순수 Go sqlite 드라이버 등록 (이름: "sqlite")

	"github.com/yoophi/refine-cms/backend/internal/config"
	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

const (
	DriverSQLite   = "sqlite"
	DriverPostgres = "postgres"
)

// NewDB 는 설정에 따라 sqlite 또는 postgres 연결을 생성한다.
// 드라이버 이름은 *Repository 들이 dialect 분기에 사용하므로 함께 반환한다.
func NewDB(cfg config.DBConfig) (*sqlx.DB, string, error) {
	switch cfg.Driver {
	case DriverPostgres:
		if cfg.DSN == "" {
			return nil, "", fmt.Errorf("postgres 드라이버에는 DB_DSN 이 필요합니다")
		}
		db, err := sqlx.Connect(DriverPostgres, cfg.DSN)
		if err != nil {
			return nil, "", fmt.Errorf("postgres 연결 실패: %w", err)
		}
		return db, DriverPostgres, nil
	case DriverSQLite, "":
		// sqlx 가 modernc 의 "sqlite" 드라이버를 ? 플레이스홀더로 인식하도록 명시 등록.
		sqlx.BindDriver(DriverSQLite, sqlx.QUESTION)
		dsn := cfg.SQLitePath + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
		db, err := sqlx.Connect(DriverSQLite, dsn)
		if err != nil {
			return nil, "", fmt.Errorf("sqlite 연결 실패: %w", err)
		}
		return db, DriverSQLite, nil
	default:
		return nil, "", fmt.Errorf("지원하지 않는 DB 드라이버입니다: %s", cfg.Driver)
	}
}

// ext 는 *sqlx.DB 와 *sqlx.Tx 가 공통으로 만족하는 실행 인터페이스이다.
// 이를 통해 헬퍼가 트랜잭션 안/밖 모두에서 동작한다.
type ext interface {
	Rebind(string) string
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
}

// insertReturningID 는 dialect 차이를 흡수해 INSERT 후 새 PK 를 반환한다.
//   - postgres: "RETURNING id" 사용
//   - sqlite:   LastInsertId() 사용
func insertReturningID(ctx context.Context, e ext, driver, query string, args ...any) (int64, error) {
	if driver == DriverPostgres {
		var id int64
		q := e.Rebind(query + " RETURNING id")
		if err := e.QueryRowxContext(ctx, q, args...).Scan(&id); err != nil {
			return 0, mapError(err)
		}
		return id, nil
	}
	res, err := e.ExecContext(ctx, e.Rebind(query), args...)
	if err != nil {
		return 0, mapError(err)
	}
	return res.LastInsertId()
}

// affectedOrNotFound 는 UPDATE/DELETE 결과의 영향 행이 0이면 ErrNotFound 를 반환한다.
func affectedOrNotFound(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// mapError 는 드라이버별 에러를 도메인 센티넬 에러로 변환한다.
func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "unique constraint failed"), // sqlite
		strings.Contains(msg, "duplicate key value"), // postgres
		strings.Contains(msg, "sqlstate 23505"):
		return domain.ErrConflict
	}
	return err
}
