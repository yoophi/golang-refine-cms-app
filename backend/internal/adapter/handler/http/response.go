package http

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// Response 는 에러 응답 표준 봉투이다. (성공 응답은 리소스 객체/`{"data": ...}` 를 직접 반환)
type Response struct {
	Code    int    `json:"code"`             // 애플리케이션 ErrorCode
	Status  int    `json:"status"`           // HTTP 상태 코드
	Message string `json:"message"`          // 사용자 노출 메시지
	Detail  string `json:"detail,omitempty"` // 4xx 에 한해 디버깅용 원본 에러 문자열
}

// respondError 는 에러를 ErrorCode 로 매핑해 gin 컨텍스트에 "등록만" 한다.
// 실제 HTTP 응답은 ErrorHandle 미들웨어가 생성한다(에러는 한 번만 처리).
// 호출 측은 이 함수 호출 후 반드시 return 한다.
func respondError(c *gin.Context, err error) {
	_ = c.Error(fromDomain(err))
}

// respondBadRequest 는 요청 바인딩/파싱 실패 등 도메인 이전 단계의 입력 오류를 등록한다.
func respondBadRequest(c *gin.Context, err error) {
	_ = c.Error(wrapGinError(err, ErrBadParamInput))
}

// parseIDParam 은 경로 파라미터 :id 를 uint 로 파싱한다.
func parseIDParam(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || v == 0 {
		respondBadRequest(c, domain.ErrInvalidInput)
		return 0, false
	}
	return uint(v), true
}
