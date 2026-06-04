package http

import (
	"time"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// ---- 요청 DTO ----

type createCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

type updateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

type createTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type updateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type createPostRequest struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Excerpt    string `json:"excerpt"`
	Content    string `json:"content"`
	Status     string `json:"status"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

type updatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Excerpt    string `json:"excerpt"`
	Content    string `json:"content"`
	Status     string `json:"status"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

type createCommentRequest struct {
	PostID      uint   `json:"post_id" binding:"required"`
	ParentID    *uint  `json:"parent_id"`
	AuthorName  string `json:"author_name" binding:"required"`
	AuthorEmail string `json:"author_email"`
	Content     string `json:"content" binding:"required"`
}

type updateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	Status  string `json:"status"`
}

// ---- 응답 DTO + 매핑 ----

type categoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newCategoryResponse(c *domain.Category) categoryResponse {
	return categoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		ParentID:    c.ParentID,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type tagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newTagResponse(t *domain.Tag) tagResponse {
	return tagResponse{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type postResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Excerpt     string        `json:"excerpt"`
	Content     string        `json:"content"`
	Status      string        `json:"status"`
	CategoryID  *uint         `json:"category_id"`
	Tags        []tagResponse `json:"tags"`
	PublishedAt *time.Time    `json:"published_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func newPostResponse(p *domain.Post) postResponse {
	tags := make([]tagResponse, 0, len(p.Tags))
	for i := range p.Tags {
		tags = append(tags, newTagResponse(&p.Tags[i]))
	}
	return postResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Excerpt:     p.Excerpt,
		Content:     p.Content,
		Status:      string(p.Status),
		CategoryID:  p.CategoryID,
		Tags:        tags,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type commentResponse struct {
	ID          uint      `json:"id"`
	PostID      uint      `json:"post_id"`
	ParentID    *uint     `json:"parent_id"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newCommentResponse(c *domain.Comment) commentResponse {
	return commentResponse{
		ID:          c.ID,
		PostID:      c.PostID,
		ParentID:    c.ParentID,
		AuthorName:  c.AuthorName,
		AuthorEmail: c.AuthorEmail,
		Content:     c.Content,
		Status:      string(c.Status),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
