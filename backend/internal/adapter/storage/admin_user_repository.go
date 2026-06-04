package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// AdminUserRepository 는 port.AdminUserRepository 의 sqlx 구현이다.
type AdminUserRepository struct {
	db     *sqlx.DB
	driver string
}

var _ port.AdminUserRepository = (*AdminUserRepository)(nil)

func NewAdminUserRepository(db *sqlx.DB, driver string) *AdminUserRepository {
	return &AdminUserRepository{db: db, driver: driver}
}

type adminUserRow struct {
	ID           uint      `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Role         string    `db:"role"`
	PasswordHash string    `db:"password_hash"`
	Avatar       string    `db:"avatar"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (r adminUserRow) toDomain() domain.AdminUser {
	return domain.AdminUser{
		ID:           r.ID,
		Name:         r.Name,
		Email:        r.Email,
		Role:         domain.Role(r.Role),
		PasswordHash: r.PasswordHash,
		Avatar:       r.Avatar,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

const adminUserCols = "id, name, email, role, password_hash, avatar, created_at, updated_at"

func (r *AdminUserRepository) GetByEmail(ctx context.Context, email string) (*domain.AdminUser, error) {
	var row adminUserRow
	const q = `SELECT ` + adminUserCols + ` FROM admin_users WHERE email = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), email); err != nil {
		return nil, errors.Wrap(mapError(err), "관리자 조회(email)")
	}
	u := row.toDomain()
	return &u, nil
}

func (r *AdminUserRepository) GetByID(ctx context.Context, id uint) (*domain.AdminUser, error) {
	var row adminUserRow
	const q = `SELECT ` + adminUserCols + ` FROM admin_users WHERE id = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), id); err != nil {
		return nil, errors.Wrap(mapError(err), "관리자 조회(id)")
	}
	u := row.toDomain()
	return &u, nil
}

// Count 는 관리자 수를 반환한다(시드 여부 판단용).
func (r *AdminUserRepository) Count(ctx context.Context) (int, error) {
	var n int
	if err := r.db.GetContext(ctx, &n, "SELECT COUNT(*) FROM admin_users"); err != nil {
		return 0, errors.Wrap(mapError(err), "관리자 개수 조회")
	}
	return n, nil
}

// Create 는 관리자 계정을 저장한다(시드/관리 목적).
func (r *AdminUserRepository) Create(ctx context.Context, u *domain.AdminUser) error {
	now := time.Now()
	u.CreatedAt, u.UpdatedAt = now, now
	const q = `INSERT INTO admin_users (name, email, role, password_hash, avatar, created_at, updated_at)
	           VALUES (?, ?, ?, ?, ?, ?, ?)`
	id, err := insertReturningID(ctx, r.db, r.driver, q,
		u.Name, u.Email, string(u.Role), u.PasswordHash, u.Avatar, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "관리자 생성")
	}
	u.ID = uint(id)
	return nil
}
