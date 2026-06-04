package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// Handlers 는 관리자 API 라우트 구성에 필요한 핸들러를 모은다.
type Handlers struct {
	Auth     *AuthHandler
	Post     *PostHandler
	Category *CategoryHandler
	Tag      *TagHandler
	Comment  *CommentHandler
}

// RegisterRoutes 는 주어진 그룹(예: /admin/api/v1)에 관리자 라우트를 등록한다.
//
// 보안 구조:
//   - ErrorHandle: 에러를 `{message}` 로 렌더링(이 그룹 범위).
//   - /auth/login, /auth/refresh 는 공개(토큰 불필요).
//   - 그 외 전부 Authenticate(액세스 토큰 검증) 적용.
//   - 각 CRUD 라우트는 RequirePermission 으로 `resource:action` 권한을 강제(401/403).
func RegisterRoutes(rg *gin.RouterGroup, auth port.AuthService, h Handlers) {
	rg.Use(ErrorHandle())

	// 공개 엔드포인트
	rg.POST("/auth/login", h.Auth.login)
	rg.POST("/auth/refresh", h.Auth.refresh)

	// 보호 엔드포인트(액세스 토큰 필요)
	sec := rg.Group("")
	sec.Use(Authenticate(auth))

	sec.GET("/auth/me", h.Auth.me)
	sec.POST("/auth/logout", h.Auth.logout)

	h.Post.register(sec)
	h.Category.register(sec)
	h.Tag.register(sec)
	h.Comment.register(sec)
}
