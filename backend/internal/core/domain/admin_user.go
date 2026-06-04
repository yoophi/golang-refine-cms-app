package domain

import "time"

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

// 역할별 권한은 불변이므로 패키지 로딩 시 한 번만 계산해 둔다(요청마다 재생성 방지).
// rolePermissionList: 응답용 펼쳐진 목록, rolePermissionSet: Can() O(1) 조회용 집합.
var (
	rolePermissionList = map[Role][]string{
		RoleSuperadmin: expandPermissions(adminResources, adminActions),
		RoleEditor:     append(expandPermissions(adminResources, []string{"list", "show", "create", "edit"}), "comments:delete"),
		RoleViewer:     expandPermissions(adminResources, []string{"list", "show"}),
	}
	rolePermissionSet = buildPermissionSets(rolePermissionList)
)

func expandPermissions(resources, actions []string) []string {
	perms := make([]string, 0, len(resources)*len(actions))
	for _, res := range resources {
		for _, a := range actions {
			perms = append(perms, res+":"+a)
		}
	}
	return perms
}

func buildPermissionSets(m map[Role][]string) map[Role]map[string]struct{} {
	out := make(map[Role]map[string]struct{}, len(m))
	for r, perms := range m {
		set := make(map[string]struct{}, len(perms))
		for _, p := range perms {
			set[p] = struct{}{}
		}
		out[r] = set
	}
	return out
}

// Permissions 는 역할에 펼쳐진 `resource:action` 권한 목록을 반환한다(사전 계산된 슬라이스).
func (r Role) Permissions() []string { return rolePermissionList[r] }

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
	_, ok := rolePermissionSet[u.Role][resource+":"+action]
	return ok
}
