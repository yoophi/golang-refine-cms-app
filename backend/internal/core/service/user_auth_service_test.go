package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
	"github.com/yoophi/refine-cms/backend/internal/core/service"
)

func newUserAuth() (port.UserAuthService, *fakeUserRepo) {
	repo := newFakeUserRepo()
	return service.NewUserAuthService(repo, fakeTokens{}, fakeHasher{}), repo
}

func TestUserAuth_RegisterAndLogin(t *testing.T) {
	svc, _ := newUserAuth()
	ctx := context.Background()

	res, err := svc.Register(ctx, port.RegisterUserInput{Email: "User@Example.com ", Password: "pw123456", Name: "홍길동"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if res.User.Email != "user@example.com" { // 정규화(소문자/trim) 확인
		t.Fatalf("email normalize: %q", res.User.Email)
	}
	if res.AccessToken == "" || res.User.ID == 0 {
		t.Fatalf("토큰/ID 누락: %+v", res)
	}

	// 중복 가입 → ErrConflict
	if _, err := svc.Register(ctx, port.RegisterUserInput{Email: "user@example.com", Password: "x12345", Name: "B"}); !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("dup register: want ErrConflict, got %v", err)
	}

	// 로그인 성공/실패
	if _, err := svc.Login(ctx, "user@example.com", "pw123456"); err != nil {
		t.Fatalf("login ok: %v", err)
	}
	if _, err := svc.Login(ctx, "user@example.com", "wrong"); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("login wrong pw: want ErrUnauthorized, got %v", err)
	}
	if _, err := svc.Login(ctx, "nobody@example.com", "pw123456"); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("login no account: want ErrUnauthorized, got %v", err)
	}
}

func TestUserAuth_Identity(t *testing.T) {
	svc, _ := newUserAuth()
	ctx := context.Background()
	res, _ := svc.Register(ctx, port.RegisterUserInput{Email: "a@example.com", Password: "pw123456", Name: "A"})

	u, err := svc.Identity(ctx, res.AccessToken)
	if err != nil || u.ID != res.User.ID {
		t.Fatalf("identity: %v, u=%+v", err, u)
	}
	if _, err := svc.Identity(ctx, "garbage"); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("identity bad token: want ErrUnauthorized, got %v", err)
	}
}

func TestUserAuth_UpdateMe_PasswordRequiresCurrent(t *testing.T) {
	svc, _ := newUserAuth()
	ctx := context.Background()
	res, _ := svc.Register(ctx, port.RegisterUserInput{Email: "a@example.com", Password: "old12345", Name: "A"})
	id := res.User.ID

	newName := "변경됨"
	if u, err := svc.UpdateMe(ctx, id, port.UpdateMeInput{Name: &newName}); err != nil || u.Name != "변경됨" {
		t.Fatalf("update name: %v, u=%+v", err, u)
	}

	pw := "new12345"
	wrong := "nope"
	// currentPassword 누락 → 401
	if _, err := svc.UpdateMe(ctx, id, port.UpdateMeInput{Password: &pw}); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("update pw without current: want ErrUnauthorized, got %v", err)
	}
	// currentPassword 틀림 → 401
	if _, err := svc.UpdateMe(ctx, id, port.UpdateMeInput{Password: &pw, CurrentPassword: &wrong}); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("update pw wrong current: want ErrUnauthorized, got %v", err)
	}
	// 올바른 currentPassword → 성공, 새 비번으로 로그인 가능
	cur := "old12345"
	if _, err := svc.UpdateMe(ctx, id, port.UpdateMeInput{Password: &pw, CurrentPassword: &cur}); err != nil {
		t.Fatalf("update pw ok: %v", err)
	}
	if _, err := svc.Login(ctx, "a@example.com", "new12345"); err != nil {
		t.Fatalf("login with new pw: %v", err)
	}
}
