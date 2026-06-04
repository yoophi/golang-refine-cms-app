package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// ErrorHandle 은 핸들러가 c.Error() 로 등록한 에러를 표준 Response 봉투로 변환한다.
// 응답 생성은 이 미들웨어 한 곳에서만 일어난다(에러는 한 번만 처리). 로깅은 ginzap 이 담당.
func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		var ginErr *GinError
		if !errors.As(err, &ginErr) || ginErr == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
				Code:    int(ErrInternal),
				Status:  http.StatusInternalServerError,
				Message: ErrInternal.Message(),
			})
			return
		}
		resp := Response{
			Code:    int(ginErr.Code()),
			Status:  ginErr.Code().StatusCode(),
			Message: ginErr.Code().Message(),
		}
		if ginErr.Code() != ErrInternal {
			resp.Detail = ginErr.Error()
		}
		c.AbortWithStatusJSON(resp.Status, resp)
	}
}

// userIdentityKey 는 인증된 회원 신원을 gin 컨텍스트에 저장하는 키이다.
const userIdentityKey = "userIdentity"

func setUserIdentity(c *gin.Context, u *domain.User) {
	c.Set(userIdentityKey, u)
}

// userIdentityOf 는 컨텍스트의 회원 신원을 반환한다(Authenticate 미들웨어가 설정).
func userIdentityOf(c *gin.Context) (*domain.User, bool) {
	v, ok := c.Get(userIdentityKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*domain.User)
	return u, ok
}

// Authenticate 는 회원 Bearer 액세스 토큰을 검증해 신원을 컨텍스트에 적재한다.
// 토큰 누락/만료/무효(또는 관리자 토큰 등 타용도) → 401(domain.ErrUnauthorized).
func Authenticate(auth port.UserAuthService) gin.HandlerFunc {
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
		setUserIdentity(c, u)
		c.Next()
	}
}
