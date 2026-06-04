package storage

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// PostRepository 는 port.PostRepository 의 sqlx 구현이다.
type PostRepository struct {
	db     *sqlx.DB
	driver string
}

var _ port.PostRepository = (*PostRepository)(nil)

func NewPostRepository(db *sqlx.DB, driver string) *PostRepository {
	return &PostRepository{db: db, driver: driver}
}

func (r *PostRepository) Create(ctx context.Context, p *domain.Post, tagIDs []uint) error {
	now := time.Now()
	p.CreatedAt, p.UpdatedAt = now, now

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(mapError(err), "게시글 생성: 트랜잭션 시작")
	}
	defer func() { _ = tx.Rollback() }()

	const q = `INSERT INTO posts (title, slug, excerpt, content, status, category_id, published_at, created_at, updated_at)
	           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	id, err := insertReturningID(ctx, tx, r.driver, q,
		p.Title, p.Slug, p.Excerpt, p.Content, string(p.Status), p.CategoryID, p.PublishedAt, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "게시글 생성")
	}
	p.ID = uint(id)

	if err := r.replaceTags(ctx, tx, p.ID, tagIDs); err != nil {
		return errors.Wrap(err, "게시글 생성: 태그 연결")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(mapError(err), "게시글 생성: 커밋")
	}
	return nil
}

func (r *PostRepository) GetByID(ctx context.Context, id uint) (*domain.Post, error) {
	var row postRow
	const q = `SELECT id, title, slug, excerpt, content, status, category_id, published_at, created_at, updated_at
	           FROM posts WHERE id = ?`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(q), id); err != nil {
		return nil, errors.Wrap(mapError(err), "게시글 조회")
	}
	p := row.toDomain()

	tagsByPost, err := r.loadTags(ctx, []uint{id})
	if err != nil {
		return nil, errors.Wrap(err, "게시글 조회: 태그 로드")
	}
	p.Tags = tagsByPost[id]
	return &p, nil
}

func (r *PostRepository) List(ctx context.Context, f port.PostFilter) ([]domain.Post, error) {
	q := `SELECT id, title, slug, excerpt, content, status, category_id, published_at, created_at, updated_at FROM posts`
	conds := make([]string, 0, 2)
	args := make([]any, 0, 2)
	if f.Status != nil {
		conds = append(conds, "status = ?")
		args = append(args, string(*f.Status))
	}
	if f.CategoryID != nil {
		conds = append(conds, "category_id = ?")
		args = append(args, *f.CategoryID)
	}
	if len(conds) > 0 {
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += " ORDER BY id DESC"

	var rows []postRow
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(q), args...); err != nil {
		return nil, errors.Wrap(mapError(err), "게시글 목록 조회")
	}
	if len(rows) == 0 {
		return []domain.Post{}, nil
	}

	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	tagsByPost, err := r.loadTags(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "게시글 목록 조회: 태그 로드")
	}

	out := make([]domain.Post, 0, len(rows))
	for _, row := range rows {
		p := row.toDomain()
		p.Tags = tagsByPost[row.ID]
		out = append(out, p)
	}
	return out, nil
}

func (r *PostRepository) Update(ctx context.Context, p *domain.Post, tagIDs []uint) error {
	p.UpdatedAt = time.Now()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(mapError(err), "게시글 수정: 트랜잭션 시작")
	}
	defer func() { _ = tx.Rollback() }()

	const q = `UPDATE posts
	           SET title = ?, slug = ?, excerpt = ?, content = ?, status = ?, category_id = ?, published_at = ?, updated_at = ?
	           WHERE id = ?`
	res, err := tx.ExecContext(ctx, tx.Rebind(q),
		p.Title, p.Slug, p.Excerpt, p.Content, string(p.Status), p.CategoryID, p.PublishedAt, p.UpdatedAt, p.ID)
	if err != nil {
		return errors.Wrap(mapError(err), "게시글 수정")
	}
	if err := affectedOrNotFound(res); err != nil {
		return errors.Wrap(err, "게시글 수정")
	}

	// tagIDs 가 nil 이면 태그 연결을 변경하지 않는다(빈 슬라이스는 전체 해제).
	if tagIDs != nil {
		if err := r.replaceTags(ctx, tx, p.ID, tagIDs); err != nil {
			return errors.Wrap(err, "게시글 수정: 태그 연결")
		}
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(mapError(err), "게시글 수정: 커밋")
	}
	return nil
}

func (r *PostRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, r.db.Rebind(`DELETE FROM posts WHERE id = ?`), id)
	if err != nil {
		return errors.Wrap(mapError(err), "게시글 삭제")
	}
	return errors.Wrap(affectedOrNotFound(res), "게시글 삭제")
}

// replaceTags 는 게시글의 태그 연결을 주어진 tagIDs 로 교체한다.
func (r *PostRepository) replaceTags(ctx context.Context, tx *sqlx.Tx, postID uint, tagIDs []uint) error {
	if _, err := tx.ExecContext(ctx, tx.Rebind(`DELETE FROM post_tags WHERE post_id = ?`), postID); err != nil {
		return errors.Wrap(mapError(err), "기존 태그 연결 삭제")
	}
	for _, tagID := range tagIDs {
		if _, err := tx.ExecContext(ctx,
			tx.Rebind(`INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)`), postID, tagID); err != nil {
			return errors.Wrapf(mapError(err), "태그 연결(tag_id=%d)", tagID)
		}
	}
	return nil
}

// loadTags 는 여러 게시글의 태그를 한 번의 쿼리로 조회해 postID 별로 묶는다(N+1 회피).
func (r *PostRepository) loadTags(ctx context.Context, postIDs []uint) (map[uint][]domain.Tag, error) {
	result := make(map[uint][]domain.Tag, len(postIDs))
	if len(postIDs) == 0 {
		return result, nil
	}
	const base = `SELECT pt.post_id AS post_id, t.id AS id, t.name AS name, t.slug AS slug,
	                     t.created_at AS created_at, t.updated_at AS updated_at
	              FROM post_tags pt
	              JOIN tags t ON t.id = pt.tag_id
	              WHERE pt.post_id IN (?)
	              ORDER BY t.id`
	query, args, err := sqlx.In(base, postIDs)
	if err != nil {
		return nil, errors.Wrap(err, "태그 IN 절 생성")
	}
	var rows []postTagRow
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		return nil, errors.Wrap(mapError(err), "태그 조인 조회")
	}
	for _, row := range rows {
		result[row.PostID] = append(result[row.PostID], row.tagRow.toDomain())
	}
	return result, nil
}
