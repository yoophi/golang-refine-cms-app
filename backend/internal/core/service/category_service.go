package service

import (
	"context"
	"strings"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

type categoryService struct {
	repo port.CategoryRepository
}

// NewCategoryService 는 CategoryService 구현을 생성한다.
func NewCategoryService(repo port.CategoryRepository) port.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, in port.CreateCategoryInput) (*domain.Category, error) {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, domain.ErrInvalidInput
	}
	c := &domain.Category{
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		ParentID:    in.ParentID,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *categoryService) Get(ctx context.Context, id uint) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *categoryService) List(ctx context.Context) ([]domain.Category, error) {
	return s.repo.List(ctx)
}

func (s *categoryService) Query(ctx context.Context, q port.ListQuery) ([]domain.Category, int, error) {
	return s.repo.Query(ctx, q)
}

func (s *categoryService) Update(ctx context.Context, id uint, in port.UpdateCategoryInput) (*domain.Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, domain.ErrInvalidInput
	}
	c.Name = in.Name
	c.Slug = in.Slug
	c.Description = in.Description
	c.ParentID = in.ParentID
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
