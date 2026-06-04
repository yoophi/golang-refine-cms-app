package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ErrorHandle 은 핸들러가 c.Error() 로 등록한 에러를 표준 HTTP 응답으로 변환한다.
//
// 핵심 원칙(에러는 딱 한 번만 처리):
//   - 핸들러는 직접 c.JSON 으로 에러 응답을 만들지 않고 c.Error(...) 로 "등록만" 한다.
//   - 에러 → HTTP 응답 변환은 오직 이 미들웨어 한 곳에서 일어난다.
//   - 에러 로깅은 별도(ginzap)에서 c.Errors 를 읽어 처리하므로 로깅/응답 책임이 분리된다.
func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		var ginErr *GinError
		if !errors.As(err, &ginErr) || ginErr == nil {
			// GinError 로 분류되지 않은 에러는 내부 오류로 처리하고 상세를 노출하지 않는다.
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
		// 내부 오류(500)는 원본 상세를 클라이언트에 노출하지 않는다.
		// 분류된 에러(4xx)만 디버깅 편의를 위해 detail 을 포함한다.
		if ginErr.Code() != ErrInternal {
			resp.Detail = ginErr.Error()
		}
		c.AbortWithStatusJSON(resp.Status, resp)
	}
}
