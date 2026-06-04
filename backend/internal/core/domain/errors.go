package domain

// ConstantError 는 const 로 선언 가능한 불변 sentinel 에러 타입이다.
// (string 기반이므로 컴파일 타임 상수가 되어 런타임에 변경/치환될 수 없다.)
// pos-connector 의 도메인 에러 컨벤션과 동일한 패턴이다.
type ConstantError string

func (e ConstantError) Error() string { return string(e) }

// 도메인 전반에서 공유하는 sentinel 에러. 식별은 errors.Is 로 한다.
// 이 값들은 어댑터 계층(HTTP)이 ErrorCode 로 매핑하므로 도메인 공개 API 의 일부이다.
// 에러를 전파할 때는 github.com/pkg/errors 의 Wrap 으로 맥락을 덧붙이되,
// errors.Is 체인은 보존되므로 깊이 래핑돼도 아래 값과의 매칭이 동작한다.
const (
	// ErrNotFound 는 요청한 리소스가 존재하지 않을 때.
	ErrNotFound = ConstantError("요청한 리소스를 찾을 수 없습니다")
	// ErrConflict 는 유니크 제약 위반 등 충돌이 발생했을 때.
	ErrConflict = ConstantError("리소스가 이미 존재합니다")
	// ErrInvalidInput 은 입력값 검증에 실패했을 때.
	ErrInvalidInput = ConstantError("입력값이 올바르지 않습니다")
)
