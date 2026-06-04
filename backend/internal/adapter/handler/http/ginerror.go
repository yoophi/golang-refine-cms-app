package http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// ErrorCode 는 애플리케이션 에러 분류이며 HTTP 상태 코드로 매핑된다.
// (pos-connector 의 gin/errors 컨벤션과 동일한 구조.)
type ErrorCode int

const (
	ErrInternal ErrorCode = iota
	ErrNotFound
	ErrBadParamInput
	ErrConflict
	ErrUnauthorized
	ErrForbidden
)

// StatusCode 는 ErrorCode 를 HTTP 상태 코드로 변환한다. (상태 매핑은 이 한 곳에 집중)
func (c ErrorCode) StatusCode() int {
	switch c {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadParamInput:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// Message 는 사용자에게 노출할 안전한 메시지이다(내부 에러 상세를 흘리지 않는다).
func (c ErrorCode) Message() string {
	switch c {
	case ErrNotFound:
		return "요청한 리소스를 찾을 수 없습니다"
	case ErrBadParamInput:
		return "입력값이 올바르지 않습니다"
	case ErrConflict:
		return "리소스가 이미 존재합니다"
	case ErrUnauthorized:
		return "인증이 필요합니다"
	case ErrForbidden:
		return "권한이 없습니다"
	default:
		return "내부 서버 오류"
	}
}

// GinError 는 원본 에러에 응답용 ErrorCode 를 부착한 래퍼이다.
// errors.As(err, &*GinError) 로 미들웨어에서 풀어낸다.
type GinError struct {
	code ErrorCode
	err  error
}

func (e *GinError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.code.Message()
}

// Unwrap 으로 원본 에러 체인을 노출해 errors.Is/As 가 통과하도록 한다.
func (e *GinError) Unwrap() error { return e.err }

func (e *GinError) Code() ErrorCode { return e.code }

// wrapGinError 는 임의의 에러에 ErrorCode 를 부착한다.
func wrapGinError(err error, code ErrorCode) *GinError {
	return &GinError{code: code, err: err}
}

// fromDomain 은 도메인 sentinel 에러(errors.Is 로 식별)를 ErrorCode 로 매핑한다.
// 알 수 없는 에러는 ErrInternal(500)로 처리한다.
func fromDomain(err error) *GinError {
	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		return wrapGinError(err, ErrUnauthorized)
	case errors.Is(err, domain.ErrForbidden):
		return wrapGinError(err, ErrForbidden)
	case errors.Is(err, domain.ErrNotFound):
		return wrapGinError(err, ErrNotFound)
	case errors.Is(err, domain.ErrConflict):
		return wrapGinError(err, ErrConflict)
	case errors.Is(err, domain.ErrInvalidInput):
		return wrapGinError(err, ErrBadParamInput)
	default:
		return wrapGinError(err, ErrInternal)
	}
}
