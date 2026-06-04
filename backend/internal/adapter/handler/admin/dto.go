package admin

import (
	"time"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// 관리자 API DTO: refine simple-rest 규약(camelCase) 으로 직렬화한다.
// 사용자(public) API 가 snake_case 인 것과 분리된 별도 표현이다.

// ---- 응답 DTO ----

type categoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parentId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	CategoryID  *uint         `json:"categoryId"`
	TagIDs      []uint        `json:"tagIds"`
	Tags        []tagResponse `json:"tags"`
	PublishedAt *time.Time    `json:"publishedAt"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

func newPostResponse(p *domain.Post) postResponse {
	tagIDs := make([]uint, 0, len(p.Tags))
	tags := make([]tagResponse, 0, len(p.Tags))
	for i := range p.Tags {
		tagIDs = append(tagIDs, p.Tags[i].ID)
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
		TagIDs:      tagIDs,
		Tags:        tags,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type commentResponse struct {
	ID          uint      `json:"id"`
	PostID      uint      `json:"postId"`
	ParentID    *uint     `json:"parentId"`
	AuthorName  string    `json:"authorName"`
	AuthorEmail string    `json:"authorEmail"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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

// ---- 생성(POST) 요청 DTO ----

type categoryCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parentId"`
}

type tagCreateRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type postCreateRequest struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Excerpt    string `json:"excerpt"`
	Content    string `json:"content"`
	Status     string `json:"status"`
	CategoryID *uint  `json:"categoryId"`
	TagIDs     []uint `json:"tagIds"`
}

type commentCreateRequest struct {
	PostID      uint   `json:"postId" binding:"required"`
	ParentID    *uint  `json:"parentId"`
	AuthorName  string `json:"authorName" binding:"required"`
	AuthorEmail string `json:"authorEmail"`
	Content     string `json:"content" binding:"required"`
	Status      string `json:"status"`
}
