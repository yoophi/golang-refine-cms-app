package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CommentRepository 는 port.CommentRepository 의 sqlx 구현이다.
type CommentRepository struct {
	db     *sqlx.DB
	driver string
}

var _ port.CommentRepository = (*CommentRepository)(nil)

func NewCommentRepository(db *sqlx.DB, driver string) *CommentRepository {
	return &CommentRepository{db: db, driver: driver}
}

func (r *CommentRepository) Create(ctx context.Context, c *domain.Comment) error {
	now := time.Now()
	c.CreatedAt, c.UpdatedAt = now, now
	const q = `INSERT INTO comments (post_id, parent_id, author_name, author_email, content, status, created_at, updated_at)
	           VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	id, err := insertReturningID(ctx, r.db, r.driver, q,
		c.PostID, c.ParentID, c.AuthorName, c.AuthorEmail, c.Content, string(c.Status), c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}
	c.ID = uint(id)
	return nil
}

func (r *CommentRepository) GetByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var row commentRow
	const q = `SELECT id, post_id, parent_id, author_name, author_email, content, status, created_at, updated_at
	           FROM comments WHERE id = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), id); err != nil {
		return nil, mapError(err)
	}
	d := row.toDomain()
	return &d, nil
}

func (r *CommentRepository) ListByPost(ctx context.Context, postID uint) ([]domain.Comment, error) {
	var rows []commentRow
	const q = `SELECT id, post_id, parent_id, author_name, author_email, content, status, created_at, updated_at
	           FROM comments WHERE post_id = ? ORDER BY id`
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(q), postID); err != nil {
		return nil, mapError(err)
	}
	out := make([]domain.Comment, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.toDomain())
	}
	return out, nil
}

func (r *CommentRepository) Update(ctx context.Context, c *domain.Comment) error {
	c.UpdatedAt = time.Now()
	const q = `UPDATE comments SET content = ?, status = ?, updated_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, r.db.Rebind(q), c.Content, string(c.Status), c.UpdatedAt, c.ID)
	if err != nil {
		return mapError(err)
	}
	return affectedOrNotFound(res)
}

func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, r.db.Rebind(`DELETE FROM comments WHERE id = ?`), id)
	if err != nil {
		return mapError(err)
	}
	return affectedOrNotFound(res)
}
