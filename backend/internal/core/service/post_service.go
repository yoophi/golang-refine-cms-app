package service

import (
	"context"
	"strings"
	"time"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

type postService struct {
	repo port.PostRepository
}

// NewPostService 는 PostService 구현을 생성한다.
func NewPostService(repo port.PostRepository) port.PostService {
	return &postService{repo: repo}
}

func (s *postService) Create(ctx context.Context, in port.CreatePostInput) (*domain.Post, error) {
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, domain.ErrInvalidInput
	}

	status := in.Status
	if status == "" {
		status = domain.PostStatusDraft
	}
	if !status.Valid() {
		return nil, domain.ErrInvalidInput
	}

	p := &domain.Post{
		Title:      in.Title,
		Slug:       in.Slug,
		Excerpt:    in.Excerpt,
		Content:    in.Content,
		Status:     status,
		CategoryID: in.CategoryID,
	}
	// 발행 상태로 생성되면 발행 시각을 기록한다.
	if status == domain.PostStatusPublished {
		now := time.Now()
		p.PublishedAt = &now
	}

	if err := s.repo.Create(ctx, p, in.TagIDs); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, p.ID)
}

func (s *postService) Get(ctx context.Context, id uint) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *postService) List(ctx context.Context, f port.PostFilter) ([]domain.Post, error) {
	return s.repo.List(ctx, f)
}

func (s *postService) Update(ctx context.Context, id uint, in port.UpdatePostInput) (*domain.Post, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, domain.ErrInvalidInput
	}

	status := in.Status
	if status == "" {
		status = p.Status
	}
	if !status.Valid() {
		return nil, domain.ErrInvalidInput
	}

	// draft -> published 로 처음 전환되는 시점에 발행 시각을 기록한다.
	if status == domain.PostStatusPublished && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}

	p.Title = in.Title
	p.Slug = in.Slug
	p.Excerpt = in.Excerpt
	p.Content = in.Content
	p.Status = status
	p.CategoryID = in.CategoryID

	if err := s.repo.Update(ctx, p, in.TagIDs); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, p.ID)
}

func (s *postService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
