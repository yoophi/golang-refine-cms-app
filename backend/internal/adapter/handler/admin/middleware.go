package admin

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CORS 는 공용 httpmw 패키지로 이동했다(공개/관리자 공통). 여기서는 관리자 인증/인가만 다룬다.

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
