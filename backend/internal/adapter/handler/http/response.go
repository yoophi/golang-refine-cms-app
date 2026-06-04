package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// errorResponse 는 표준 에러 응답 바디이다.
type errorResponse struct {
	Error string `json:"error"`
}

// respondError 는 도메인 에러를 적절한 HTTP 상태 코드로 매핑해 응답한다.
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrConflict):
		c.JSON(http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "내부 서버 오류"})
	}
}

// parseIDParam 은 경로 파라미터 :id 를 uint 로 파싱한다.
func parseIDParam(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || v == 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "잘못된 ID 입니다"})
		return 0, false
	}
	return uint(v), true
}
