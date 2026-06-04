package service_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// ---- fake PasswordHasher ----

type fakeHasher struct{}

func (fakeHasher) Hash(plain string) (string, error) { return "h:" + plain, nil }
func (fakeHasher) Compare(hash, plain string) error {
	if hash != "h:"+plain {
		return domain.ErrUnauthorized
	}
	return nil
}

// ---- fake TokenManager (토큰에 subject 를 인코딩) ----

type fakeTokens struct{}

func (fakeTokens) Issue(subject uint, role string) (string, string, error) {
	return fmt.Sprintf("a-%d", subject), fmt.Sprintf("r-%d", subject), nil
}
func (fakeTokens) ParseAccess(token string) (uint, string, error) {
	return parseFakeToken(token, "a-")
}
func (fakeTokens) ParseRefresh(token string) (uint, error) {
	id, _, err := parseFakeToken(token, "r-")
	return id, err
}
func parseFakeToken(token, prefix string) (uint, string, error) {
	if !strings.HasPrefix(token, prefix) {
		return 0, "", domain.ErrUnauthorized
	}
	n, err := strconv.ParseUint(strings.TrimPrefix(token, prefix), 10, 64)
	if err != nil {
		return 0, "", domain.ErrUnauthorized
	}
	return uint(n), "", nil
}

// ---- fake UserRepository ----

type fakeUserRepo struct {
	items  map[uint]*domain.User
	nextID uint
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{items: map[uint]*domain.User{}, nextID: 1}
}

func (f *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	for _, u := range f.items {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeUserRepo) GetByID(_ context.Context, id uint) (*domain.User, error) {
	u, ok := f.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (f *fakeUserRepo) Create(_ context.Context, u *domain.User) error {
	u.ID = f.nextID
	f.nextID++
	cp := *u
	f.items[u.ID] = &cp
	return nil
}

func (f *fakeUserRepo) Update(_ context.Context, u *domain.User) error {
	if _, ok := f.items[u.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *u
	f.items[u.ID] = &cp
	return nil
}

// ---- fake CommentRepository ----

type fakeCommentRepo struct {
	items  map[uint]*domain.Comment
	nextID uint
}

func newFakeCommentRepo() *fakeCommentRepo {
	return &fakeCommentRepo{items: map[uint]*domain.Comment{}, nextID: 1}
}

func (f *fakeCommentRepo) Create(_ context.Context, c *domain.Comment) error {
	c.ID = f.nextID
	f.nextID++
	cp := *c
	f.items[c.ID] = &cp
	return nil
}

func (f *fakeCommentRepo) GetByID(_ context.Context, id uint) (*domain.Comment, error) {
	c, ok := f.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *c
	return &cp, nil
}

func (f *fakeCommentRepo) ListByPost(_ context.Context, postID uint) ([]domain.Comment, error) {
	var out []domain.Comment
	for _, c := range f.items {
		if c.PostID == postID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (f *fakeCommentRepo) Query(_ context.Context, _ port.ListQuery) ([]domain.Comment, int, error) {
	out := make([]domain.Comment, 0, len(f.items))
	for _, c := range f.items {
		out = append(out, *c)
	}
	return out, len(out), nil
}

func (f *fakeCommentRepo) Update(_ context.Context, c *domain.Comment) error {
	if _, ok := f.items[c.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *c
	f.items[c.ID] = &cp
	return nil
}

func (f *fakeCommentRepo) Delete(_ context.Context, id uint) error {
	if _, ok := f.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.items, id)
	return nil
}

// ---- fake PostRepository (댓글 서비스가 게시글 존재 확인에만 사용) ----

type fakePostRepo struct {
	exists map[uint]bool
}

func newFakePostRepo(ids ...uint) *fakePostRepo {
	m := map[uint]bool{}
	for _, id := range ids {
		m[id] = true
	}
	return &fakePostRepo{exists: m}
}

func (f *fakePostRepo) Create(context.Context, *domain.Post, []uint) error { return nil }
func (f *fakePostRepo) GetByID(_ context.Context, id uint) (*domain.Post, error) {
	if !f.exists[id] {
		return nil, domain.ErrNotFound
	}
	return &domain.Post{ID: id}, nil
}
func (f *fakePostRepo) List(context.Context, port.PostFilter) ([]domain.Post, error) { return nil, nil }
func (f *fakePostRepo) Query(context.Context, port.ListQuery) ([]domain.Post, int, error) {
	return nil, 0, nil
}
func (f *fakePostRepo) Update(context.Context, *domain.Post, []uint) error { return nil }
func (f *fakePostRepo) Delete(context.Context, uint) error                 { return nil }
