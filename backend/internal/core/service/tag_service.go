package service

import (
	"context"
	"strings"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

type tagService struct {
	repo port.TagRepository
}

// NewTagService 는 TagService 구현을 생성한다.
func NewTagService(repo port.TagRepository) port.TagService {
	return &tagService{repo: repo}
}

func (s *tagService) Create(ctx context.Context, in port.CreateTagInput) (*domain.Tag, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, domain.ErrInvalidInput
	}
	t := &domain.Tag{Name: in.Name, Slug: in.Slug}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *tagService) Get(ctx context.Context, id uint) (*domain.Tag, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *tagService) List(ctx context.Context) ([]domain.Tag, error) {
	return s.repo.List(ctx)
}

func (s *tagService) Update(ctx context.Context, id uint, in port.UpdateTagInput) (*domain.Tag, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, domain.ErrInvalidInput
	}
	t.Name = in.Name
	t.Slug = in.Slug
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *tagService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
