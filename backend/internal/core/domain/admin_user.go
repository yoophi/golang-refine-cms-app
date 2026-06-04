package domain

import (
	"slices"
	"time"
)

// Role 은 관리자 역할이다. BE 가 역할→권한 정의의 source of truth 이다.
type Role string

const (
	RoleSuperadmin Role = "superadmin" // 전체 권한(ACL 우회)
	RoleEditor     Role = "editor"     // 모든 리소스 list/show/create/edit + comments:delete
	RoleViewer     Role = "viewer"     // 모든 리소스 list/show
)

// Valid 는 정의된 역할인지 검증한다.
func (r Role) Valid() bool {
	switch r {
	case RoleSuperadmin, RoleEditor, RoleViewer:
		return true
	default:
		return false
	}
}

// 관리자 API 가 ACL 로 보호하는 리소스/액션.
var (
	adminResources = []string{"posts", "categories", "tags", "comments"}
	adminActions   = []string{"list", "show", "create", "edit", "delete"}
)

// Permissions 는 역할에 펼쳐진 `resource:action` 권한 목록을 반환한다.
func (r Role) Permissions() []string {
	switch r {
	case RoleSuperadmin:
		return expandPermissions(adminResources, adminActions)
	case RoleEditor:
		perms := expandPermissions(adminResources, []string{"list", "show", "create", "edit"})
		return append(perms, "comments:delete")
	case RoleViewer:
		return expandPermissions(adminResources, []string{"list", "show"})
	default:
		return []string{}
	}
}

func expandPermissions(resources, actions []string) []string {
	perms := make([]string, 0, len(resources)*len(actions))
	for _, res := range resources {
		for _, a := range actions {
			perms = append(perms, res+":"+a)
		}
	}
	return perms
}

// AdminUser 는 관리자 대시보드 사용자이다.
type AdminUser struct {
	ID           uint
	Name         string
	Email        string
	Role         Role
	PasswordHash string // 응답에 절대 노출하지 않는다
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Permissions 는 사용자의 펼쳐진 권한 목록을 반환한다(역할 기반).
func (u *AdminUser) Permissions() []string { return u.Role.Permissions() }

// Can 은 사용자가 `resource:action` 을 수행할 수 있는지 판정한다(superadmin 우회).
func (u *AdminUser) Can(resource, action string) bool {
	if u.Role == RoleSuperadmin {
		return true
	}
	return slices.Contains(u.Role.Permissions(), resource+":"+action)
}
