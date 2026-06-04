package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// 관리자 API 에러 규약: 바디는 `{ "message": "..." }` (스펙 §3.3).
// 핸들러는 c.Error 로 에러를 등록만 하고, 변환은 ErrorHandle 미들웨어 한 곳에서 한다.
// (사용자 API 의 Response 봉투와 다른 포맷이므로 별도 파이프라인.)

// apiError 는 도메인 에러가 아닌 입력 단계 오류(바인딩/파싱)에 상태코드를 부여한다.
type apiError struct {
	status int
	msg    string
}

func (e apiError) Error() string { return e.msg }

// ErrorHandle 은 등록된 마지막 에러를 `{message}` 응답으로 변환한다.
func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		var ae apiError
		if errors.As(err, &ae) {
			c.AbortWithStatusJSON(ae.status, gin.H{"message": ae.msg})
			return
		}
		status, msg := classify(err)
		c.AbortWithStatusJSON(status, gin.H{"message": msg})
	}
}

// classify 는 도메인 sentinel 에러(errors.Is, pkg/errors 래핑 통과)를 상태코드/메시지로 매핑한다.
func classify(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, "인증이 필요합니다"
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, "권한이 없습니다"
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "요청한 리소스를 찾을 수 없습니다"
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, "리소스가 이미 존재합니다"
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest, "입력값이 올바르지 않습니다"
	default:
		return http.StatusInternalServerError, "내부 서버 오류"
	}
}

// respondError 는 도메인/서비스 에러를 등록한다(상태코드는 ErrorHandle 이 결정).
func respondError(c *gin.Context, err error) {
	_ = c.Error(err)
}

// respondBadRequest 는 입력 오류를 400 으로 등록하며 상세 메시지를 노출한다.
func respondBadRequest(c *gin.Context, err error) {
	_ = c.Error(apiError{status: http.StatusBadRequest, msg: err.Error()})
}

// parseID 는 경로 파라미터 :id 를 uint 로 파싱한다.
func parseID(c *gin.Context) (uint, bool) {
	v, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || v == 0 {
		respondBadRequest(c, errors.New("잘못된 ID 입니다"))
		return 0, false
	}
	return uint(v), true
}
