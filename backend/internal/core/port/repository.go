package port

import (
	"context"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// 출력 포트(Driven Port): 코어가 영속성 계층에 요구하는 계약.
// 구현은 adapter/storage 에 위치한다.

// CategoryRepository 는 카테고리 영속성 계약이다.
type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error
	GetByID(ctx context.Context, id uint) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, id uint) error
}

// TagRepository 는 태그 영속성 계약이다.
type TagRepository interface {
	Create(ctx context.Context, t *domain.Tag) error
	GetByID(ctx context.Context, id uint) (*domain.Tag, error)
	List(ctx context.Context) ([]domain.Tag, error)
	Update(ctx context.Context, t *domain.Tag) error
	Delete(ctx context.Context, id uint) error
}

// PostFilter 는 게시글 목록 조회 시 사용하는 선택적 필터이다.
type PostFilter struct {
	Status     *domain.PostStatus
	CategoryID *uint
}

// PostRepository 는 게시글 영속성 계약이다.
type PostRepository interface {
	// Create 는 게시글을 저장한다. tagIDs 에 해당하는 기존 태그를 연결한다.
	Create(ctx context.Context, p *domain.Post, tagIDs []uint) error
	GetByID(ctx context.Context, id uint) (*domain.Post, error)
	List(ctx context.Context, f PostFilter) ([]domain.Post, error)
	// Update 는 게시글을 갱신한다. tagIDs 가 nil 이 아니면 태그 연결을 교체한다.
	Update(ctx context.Context, p *domain.Post, tagIDs []uint) error
	Delete(ctx context.Context, id uint) error
}

// CommentRepository 는 댓글 영속성 계약이다.
type CommentRepository interface {
	Create(ctx context.Context, c *domain.Comment) error
	GetByID(ctx context.Context, id uint) (*domain.Comment, error)
	ListByPost(ctx context.Context, postID uint) ([]domain.Comment, error)
	Update(ctx context.Context, c *domain.Comment) error
	Delete(ctx context.Context, id uint) error
}
