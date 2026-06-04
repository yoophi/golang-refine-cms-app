package port

import (
	"context"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// 입력 포트(Driving Port): 인바운드 어댑터(HTTP)가 호출하는 유스케이스 계약.
// 구현은 core/service 에 위치한다.

// --- Category ---

type CreateCategoryInput struct {
	Name        string
	Slug        string
	Description string
	ParentID    *uint
}

type UpdateCategoryInput struct {
	Name        string
	Slug        string
	Description string
	ParentID    *uint
}

type CategoryService interface {
	Create(ctx context.Context, in CreateCategoryInput) (*domain.Category, error)
	Get(ctx context.Context, id uint) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, id uint, in UpdateCategoryInput) (*domain.Category, error)
	Delete(ctx context.Context, id uint) error
}

// --- Tag ---

type CreateTagInput struct {
	Name string
	Slug string
}

type UpdateTagInput struct {
	Name string
	Slug string
}

type TagService interface {
	Create(ctx context.Context, in CreateTagInput) (*domain.Tag, error)
	Get(ctx context.Context, id uint) (*domain.Tag, error)
	List(ctx context.Context) ([]domain.Tag, error)
	Update(ctx context.Context, id uint, in UpdateTagInput) (*domain.Tag, error)
	Delete(ctx context.Context, id uint) error
}

// --- Post ---

type CreatePostInput struct {
	Title      string
	Slug       string
	Excerpt    string
	Content    string
	Status     domain.PostStatus
	CategoryID *uint
	TagIDs     []uint
}

type UpdatePostInput struct {
	Title      string
	Slug       string
	Excerpt    string
	Content    string
	Status     domain.PostStatus
	CategoryID *uint
	TagIDs     []uint // nil 이면 태그 연결 변경 없음
}

type PostService interface {
	Create(ctx context.Context, in CreatePostInput) (*domain.Post, error)
	Get(ctx context.Context, id uint) (*domain.Post, error)
	List(ctx context.Context, f PostFilter) ([]domain.Post, error)
	Update(ctx context.Context, id uint, in UpdatePostInput) (*domain.Post, error)
	Delete(ctx context.Context, id uint) error
}

// --- Comment ---

type CreateCommentInput struct {
	PostID      uint
	ParentID    *uint
	AuthorName  string
	AuthorEmail string
	Content     string
}

type UpdateCommentInput struct {
	Content string
	Status  domain.CommentStatus
}

type CommentService interface {
	Create(ctx context.Context, in CreateCommentInput) (*domain.Comment, error)
	Get(ctx context.Context, id uint) (*domain.Comment, error)
	ListByPost(ctx context.Context, postID uint) ([]domain.Comment, error)
	Update(ctx context.Context, id uint, in UpdateCommentInput) (*domain.Comment, error)
	Delete(ctx context.Context, id uint) error
}
