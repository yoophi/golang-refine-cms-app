package storage

import (
	"time"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// row 구조체: DB 컬럼과 매핑되는 영속성 표현. 도메인 엔티티와 분리한다.

type categoryRow struct {
	ID          uint      `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	ParentID    *uint     `db:"parent_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r categoryRow) toDomain() domain.Category {
	return domain.Category{
		ID:          r.ID,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		ParentID:    r.ParentID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type tagRow struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r tagRow) toDomain() domain.Tag {
	return domain.Tag{
		ID:        r.ID,
		Name:      r.Name,
		Slug:      r.Slug,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

type postRow struct {
	ID          uint       `db:"id"`
	Title       string     `db:"title"`
	Slug        string     `db:"slug"`
	Excerpt     string     `db:"excerpt"`
	Content     string     `db:"content"`
	Status      string     `db:"status"`
	CategoryID  *uint      `db:"category_id"`
	PublishedAt *time.Time `db:"published_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (r postRow) toDomain() domain.Post {
	return domain.Post{
		ID:          r.ID,
		Title:       r.Title,
		Slug:        r.Slug,
		Excerpt:     r.Excerpt,
		Content:     r.Content,
		Status:      domain.PostStatus(r.Status),
		CategoryID:  r.CategoryID,
		PublishedAt: r.PublishedAt,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type commentRow struct {
	ID          uint      `db:"id"`
	PostID      uint      `db:"post_id"`
	ParentID    *uint     `db:"parent_id"`
	UserID      *uint     `db:"user_id"`
	AuthorName  string    `db:"author_name"`
	AuthorEmail string    `db:"author_email"`
	Content     string    `db:"content"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r commentRow) toDomain() domain.Comment {
	return domain.Comment{
		ID:          r.ID,
		PostID:      r.PostID,
		ParentID:    r.ParentID,
		UserID:      r.UserID,
		AuthorName:  r.AuthorName,
		AuthorEmail: r.AuthorEmail,
		Content:     r.Content,
		Status:      domain.CommentStatus(r.Status),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// postTagRow 는 post_tags 조인 + tags 조회 결과를 담는다.
type postTagRow struct {
	PostID uint `db:"post_id"`
	tagRow
}
