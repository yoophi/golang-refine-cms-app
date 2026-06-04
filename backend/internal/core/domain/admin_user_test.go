package domain_test

import (
	"testing"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

func TestAdminUserCan(t *testing.T) {
	cases := []struct {
		role     domain.Role
		resource string
		action   string
		want     bool
	}{
		{domain.RoleSuperadmin, "posts", "delete", true},   // superadmin 우회
		{domain.RoleSuperadmin, "comments", "delete", true},
		{domain.RoleEditor, "posts", "edit", true},
		{domain.RoleEditor, "posts", "create", true},
		{domain.RoleEditor, "posts", "delete", false},      // editor 는 글 삭제 불가
		{domain.RoleEditor, "comments", "delete", true},    // 단, 댓글 삭제는 허용
		{domain.RoleViewer, "posts", "list", true},
		{domain.RoleViewer, "posts", "show", true},
		{domain.RoleViewer, "posts", "create", false},      // viewer 는 읽기만
		{domain.RoleViewer, "comments", "edit", false},
		{domain.Role("unknown"), "posts", "list", false},   // 미정의 역할은 전부 거부
	}
	for _, c := range cases {
		u := domain.AdminUser{Role: c.role}
		if got := u.Can(c.resource, c.action); got != c.want {
			t.Errorf("role=%s %s:%s → got %v, want %v", c.role, c.resource, c.action, got, c.want)
		}
	}
}

func TestRolePermissionsCount(t *testing.T) {
	cases := map[domain.Role]int{
		domain.RoleSuperadmin: 20, // 4 resources × 5 actions
		domain.RoleEditor:     17, // 4×4 + comments:delete
		domain.RoleViewer:     8,  // 4×2 (list/show)
		domain.Role("nope"):   0,
	}
	for role, want := range cases {
		if got := len(role.Permissions()); got != want {
			t.Errorf("role=%s permissions count = %d, want %d", role, got, want)
		}
	}
}
