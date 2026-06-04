package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// UserRepository 는 port.UserRepository 의 sqlx 구현이다.
type UserRepository struct {
	db     *sqlx.DB
	driver string
}

var _ port.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sqlx.DB, driver string) *UserRepository {
	return &UserRepository{db: db, driver: driver}
}

type userRow struct {
	ID           uint      `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	Avatar       string    `db:"avatar"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (r userRow) toDomain() domain.User {
	return domain.User{
		ID:           r.ID,
		Email:        r.Email,
		Name:         r.Name,
		PasswordHash: r.PasswordHash,
		Avatar:       r.Avatar,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

const userCols = "id, email, name, password_hash, avatar, created_at, updated_at"

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var row userRow
	const q = `SELECT ` + userCols + ` FROM users WHERE email = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), email); err != nil {
		return nil, errors.Wrap(mapError(err), "회원 조회(email)")
	}
	u := row.toDomain()
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var row userRow
	const q = `SELECT ` + userCols + ` FROM users WHERE id = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), id); err != nil {
		return nil, errors.Wrap(mapError(err), "회원 조회(id)")
	}
	u := row.toDomain()
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	now := time.Now()
	u.CreatedAt, u.UpdatedAt = now, now
	const q = `INSERT INTO users (email, name, password_hash, avatar, created_at, updated_at)
	           VALUES (?, ?, ?, ?, ?, ?)`
	id, err := insertReturningID(ctx, r.db, r.driver, q,
		u.Email, u.Name, u.PasswordHash, u.Avatar, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "회원 생성")
	}
	u.ID = uint(id)
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	u.UpdatedAt = time.Now()
	const q = `UPDATE users SET email = ?, name = ?, password_hash = ?, avatar = ?, updated_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, r.db.Rebind(q),
		u.Email, u.Name, u.PasswordHash, u.Avatar, u.UpdatedAt, u.ID)
	if err != nil {
		return errors.Wrap(mapError(err), "회원 수정")
	}
	return errors.Wrap(affectedOrNotFound(res), "회원 수정")
}
