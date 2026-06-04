package domain

import "errors"

// 도메인 전반에서 공유하는 센티넬 에러. 어댑터 계층(HTTP)에서 적절한
// 상태 코드로 매핑한다.
var (
	// ErrNotFound 는 요청한 리소스가 존재하지 않을 때 반환한다.
	ErrNotFound = errors.New("리소스를 찾을 수 없습니다")
	// ErrConflict 는 유니크 제약 위반 등 충돌이 발생했을 때 반환한다.
	ErrConflict = errors.New("리소스가 이미 존재합니다")
	// ErrInvalidInput 은 입력값 검증에 실패했을 때 반환한다.
	ErrInvalidInput = errors.New("입력값이 올바르지 않습니다")
)
