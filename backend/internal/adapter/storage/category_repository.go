package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CategoryRepository 는 port.CategoryRepository 의 sqlx 구현이다.
type CategoryRepository struct {
	db     *sqlx.DB
	driver string
}

// 컴파일 타임 인터페이스 만족 확인.
var _ port.CategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(db *sqlx.DB, driver string) *CategoryRepository {
	return &CategoryRepository{db: db, driver: driver}
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	now := time.Now()
	c.CreatedAt, c.UpdatedAt = now, now
	const q = `INSERT INTO categories (name, slug, description, parent_id, created_at, updated_at)
	           VALUES (?, ?, ?, ?, ?, ?)`
	id, err := insertReturningID(ctx, r.db, r.driver, q,
		c.Name, c.Slug, c.Description, c.ParentID, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}
	c.ID = uint(id)
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (*domain.Category, error) {
	var row categoryRow
	const q = `SELECT id, name, slug, description, parent_id, created_at, updated_at
	           FROM categories WHERE id = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), id); err != nil {
		return nil, mapError(err)
	}
	d := row.toDomain()
	return &d, nil
}

func (r *CategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	var rows []categoryRow
	const q = `SELECT id, name, slug, description, parent_id, created_at, updated_at
	           FROM categories ORDER BY id`
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(q)); err != nil {
		return nil, mapError(err)
	}
	out := make([]domain.Category, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.toDomain())
	}
	return out, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	c.UpdatedAt = time.Now()
	const q = `UPDATE categories
	           SET name = ?, slug = ?, description = ?, parent_id = ?, updated_at = ?
	           WHERE id = ?`
	res, err := r.db.ExecContext(ctx, r.db.Rebind(q),
		c.Name, c.Slug, c.Description, c.ParentID, c.UpdatedAt, c.ID)
	if err != nil {
		return mapError(err)
	}
	return affectedOrNotFound(res)
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, r.db.Rebind(`DELETE FROM categories WHERE id = ?`), id)
	if err != nil {
		return mapError(err)
	}
	return affectedOrNotFound(res)
}
