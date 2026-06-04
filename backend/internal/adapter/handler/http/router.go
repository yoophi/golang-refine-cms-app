package http

import (
	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// Handlers 는 사용자(public) API 라우트 구성에 필요한 핸들러를 모은다.
type Handlers struct {
	Category *CategoryHandler
	Tag      *TagHandler
	Post     *PostHandler
	Comment  *CommentHandler
	UserAuth *UserAuthHandler
}

// RegisterRoutes 는 주어진 라우터 그룹(예: /api/v1)에 사용자 API 라우트를 등록한다.
// 에러 응답 변환(ErrorHandle)은 이 그룹 범위에만 적용된다(관리자 API 와 분리).
// 엔진/전역 미들웨어(로깅·복구·CORS)는 bootstrap 에서 구성한다.
func RegisterRoutes(rg *gin.RouterGroup, h Handlers, userAuth port.UserAuthService, authRateLimit gin.HandlerFunc) {
	rg.Use(ErrorHandle())

	auth := Authenticate(userAuth)

	// 회원 인증: register/login/refresh 는 공개(레이트 리밋 적용), me/logout 은 보호.
	rg.POST("/auth/register", authRateLimit, h.UserAuth.register)
	rg.POST("/auth/login", authRateLimit, h.UserAuth.login)
	rg.POST("/auth/refresh", authRateLimit, h.UserAuth.refresh)
	rg.GET("/auth/me", auth, h.UserAuth.me)
	rg.PATCH("/auth/me", auth, h.UserAuth.updateMe)
	rg.POST("/auth/logout", auth, h.UserAuth.logout)

	h.Category.register(rg)
	h.Tag.register(rg)
	h.Post.register(rg)
	h.Comment.register(rg, auth) // 작성/수정/삭제만 로그인 + 소유권 강제
}
