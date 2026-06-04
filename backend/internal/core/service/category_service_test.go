package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
	"github.com/yoophi/refine-cms/backend/internal/core/service"
)

// fakeCategoryRepo 는 DB 없이 서비스 로직을 검증하기 위한 인메모리 가짜 구현이다.
// 헥사고날 구조 덕분에 출력 포트만 대체하면 코어를 격리 테스트할 수 있다.
type fakeCategoryRepo struct {
	items  map[uint]*domain.Category
	nextID uint
}

func newFakeCategoryRepo() *fakeCategoryRepo {
	return &fakeCategoryRepo{items: map[uint]*domain.Category{}, nextID: 1}
}

func (f *fakeCategoryRepo) Create(_ context.Context, c *domain.Category) error {
	c.ID = f.nextID
	f.nextID++
	f.items[c.ID] = c
	return nil
}

func (f *fakeCategoryRepo) GetByID(_ context.Context, id uint) (*domain.Category, error) {
	c, ok := f.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return c, nil
}

func (f *fakeCategoryRepo) List(_ context.Context) ([]domain.Category, error) {
	out := make([]domain.Category, 0, len(f.items))
	for _, c := range f.items {
		out = append(out, *c)
	}
	return out, nil
}

func (f *fakeCategoryRepo) Query(_ context.Context, _ port.ListQuery) ([]domain.Category, int, error) {
	out := make([]domain.Category, 0, len(f.items))
	for _, c := range f.items {
		out = append(out, *c)
	}
	return out, len(out), nil
}

func (f *fakeCategoryRepo) Update(_ context.Context, c *domain.Category) error {
	if _, ok := f.items[c.ID]; !ok {
		return domain.ErrNotFound
	}
	f.items[c.ID] = c
	return nil
}

func (f *fakeCategoryRepo) Delete(_ context.Context, id uint) error {
	if _, ok := f.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.items, id)
	return nil
}

func TestCategoryService_Create(t *testing.T) {
	svc := service.NewCategoryService(newFakeCategoryRepo())

	t.Run("정상 생성", func(t *testing.T) {
		got, err := svc.Create(context.Background(), port.CreateCategoryInput{Name: "기술", Slug: "tech"})
		if err != nil {
			t.Fatalf("예상치 못한 에러: %v", err)
		}
		if got.ID == 0 || got.Name != "기술" || got.Slug != "tech" {
			t.Fatalf("결과가 올바르지 않음: %+v", got)
		}
	})

	t.Run("이름 누락 시 ErrInvalidInput", func(t *testing.T) {
		_, err := svc.Create(context.Background(), port.CreateCategoryInput{Name: " ", Slug: "tech"})
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Fatalf("ErrInvalidInput 기대, 실제: %v", err)
		}
	})
}

func TestCategoryService_GetNotFound(t *testing.T) {
	svc := service.NewCategoryService(newFakeCategoryRepo())
	_, err := svc.Get(context.Background(), 42)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("ErrNotFound 기대, 실제: %v", err)
	}
}
