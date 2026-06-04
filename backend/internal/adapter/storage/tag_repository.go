package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// TagRepository 는 port.TagRepository 의 sqlx 구현이다.
type TagRepository struct {
	db     *sqlx.DB
	driver string
}

var _ port.TagRepository = (*TagRepository)(nil)

func NewTagRepository(db *sqlx.DB, driver string) *TagRepository {
	return &TagRepository{db: db, driver: driver}
}

func (r *TagRepository) Create(ctx context.Context, t *domain.Tag) error {
	now := time.Now()
	t.CreatedAt, t.UpdatedAt = now, now
	const q = `INSERT INTO tags (name, slug, created_at, updated_at) VALUES (?, ?, ?, ?)`
	id, err := insertReturningID(ctx, r.db, r.driver, q, t.Name, t.Slug, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "태그 생성")
	}
	t.ID = uint(id)
	return nil
}

func (r *TagRepository) GetByID(ctx context.Context, id uint) (*domain.Tag, error) {
	var row tagRow
	const q = `SELECT id, name, slug, created_at, updated_at FROM tags WHERE id = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), id); err != nil {
		return nil, errors.Wrap(mapError(err), "태그 조회")
	}
	d := row.toDomain()
	return &d, nil
}

func (r *TagRepository) List(ctx context.Context) ([]domain.Tag, error) {
	var rows []tagRow
	const q = `SELECT id, name, slug, created_at, updated_at FROM tags ORDER BY id`
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(q)); err != nil {
		return nil, errors.Wrap(mapError(err), "태그 목록 조회")
	}
	out := make([]domain.Tag, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.toDomain())
	}
	return out, nil
}

func (r *TagRepository) Update(ctx context.Context, t *domain.Tag) error {
	t.UpdatedAt = time.Now()
	const q = `UPDATE tags SET name = ?, slug = ?, updated_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, r.db.Rebind(q), t.Name, t.Slug, t.UpdatedAt, t.ID)
	if err != nil {
		return errors.Wrap(mapError(err), "태그 수정")
	}
	return errors.Wrap(affectedOrNotFound(res), "태그 수정")
}

func (r *TagRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, r.db.Rebind(`DELETE FROM tags WHERE id = ?`), id)
	if err != nil {
		return errors.Wrap(mapError(err), "태그 삭제")
	}
	return errors.Wrap(affectedOrNotFound(res), "태그 삭제")
}
