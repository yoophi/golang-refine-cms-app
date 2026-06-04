package admin

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CORS 는 대시보드(다른 Origin)에서의 호출을 허용하는 미들웨어를 만든다.
// refine 페이지네이션을 위해 X-Total-Count 헤더를 노출하고, Authorization 요청 헤더를 허용한다.
// 전역(engine.Use)으로 등록해 프리플라이트(OPTIONS, 미매칭 라우트 포함)까지 처리한다.
func CORS(origins []string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowMethods:  []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{totalCountHeader},
	}
	if containsWildcard(origins) {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = origins
	}
	return cors.New(cfg)
}

// containsWildcard 는 명시적 '*' 가 있을 때만 전체 허용한다.
// 빈 목록은 전체 허용이 아니라 '교차 출처 차단'(fail-closed)으로 둔다.
func containsWildcard(origins []string) bool {
	return slices.Contains(origins, "*")
}

// Authenticate 는 Bearer 액세스 토큰을 검증하고 관리자 신원을 컨텍스트에 저장한다.
// 토큰 누락/만료/무효 → 401(domain.ErrUnauthorized). 실패는 ErrorHandle 이 렌더링한다.
func Authenticate(auth port.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if token == "" || token == header {
			respondError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		u, err := auth.Identity(c.Request.Context(), token)
		if err != nil {
			respondError(c, err)
			c.Abort()
			return
		}
		setIdentity(c, u)
		c.Next()
	}
}

// RequirePermission 은 인증된 관리자가 `resource:action` 권한을 갖는지 검사한다.
// 권한 부족 → 403(domain.ErrForbidden). superadmin 은 우회한다(domain.AdminUser.Can).
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := identityOf(c)
		if !ok {
			respondError(c, domain.ErrUnauthorized)
			c.Abort()
			return
		}
		if !u.Can(resource, action) {
			respondError(c, domain.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
