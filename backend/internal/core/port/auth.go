package port

import (
	"context"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// --- 출력 포트(Driven) ---

// AdminUserRepository 는 관리자 사용자 영속성 계약이다.
type AdminUserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.AdminUser, error)
	GetByID(ctx context.Context, id uint) (*domain.AdminUser, error)
}

// UserRepository 는 회원(공개 사용자) 영속성 계약이다.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error
}

// TokenManager 는 JWT 발급/검증 계약이다(인프라 어댑터가 구현).
// audience(회원/관리자)별로 별도 인스턴스를 사용해 상호 토큰 사용을 차단한다.
type TokenManager interface {
	// Issue 는 subject(사용자 id)와 role 로 액세스/리프레시 토큰을 발급한다(role 없으면 빈 문자열).
	Issue(subject uint, role string) (accessToken, refreshToken string, err error)
	// ParseAccess 는 액세스 토큰을 검증하고 (subject, role) 을 반환한다(무효/만료/타용도 시 domain.ErrUnauthorized).
	ParseAccess(token string) (subject uint, role string, err error)
	// ParseRefresh 는 리프레시 토큰을 검증하고 subject 를 반환한다.
	ParseRefresh(token string) (subject uint, err error)
}

// PasswordHasher 는 비밀번호 해시/검증 계약이다.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	// Compare 는 불일치 시 domain.ErrUnauthorized 를 반환한다.
	Compare(hash, plain string) error
}

// --- 입력 포트(Driving) ---

// LoginResult 는 관리자 로그인/리프레시 결과이다.
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *domain.AdminUser
}

// AuthService 는 관리자 인증 유스케이스이다.
type AuthService interface {
	Login(ctx context.Context, email, password string) (*LoginResult, error)
	// Identity 는 액세스 토큰으로 현재 관리자 신원을 복원한다(미들웨어/`/auth/me`).
	Identity(ctx context.Context, accessToken string) (*domain.AdminUser, error)
	Refresh(ctx context.Context, refreshToken string) (*LoginResult, error)
}

// UserLoginResult 는 회원 가입/로그인/리프레시 결과이다.
type UserLoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *domain.User
}

// RegisterUserInput 은 회원 가입 입력이다.
type RegisterUserInput struct {
	Email    string
	Password string
	Name     string
}

// UpdateMeInput 은 회원 정보 부분 수정 입력이다(빈 값은 변경 없음).
type UpdateMeInput struct {
	Name            *string
	Email           *string
	Password        *string
	CurrentPassword *string
}

// UserAuthService 는 회원(공개) 인증 유스케이스이다.
type UserAuthService interface {
	Register(ctx context.Context, in RegisterUserInput) (*UserLoginResult, error)
	Login(ctx context.Context, email, password string) (*UserLoginResult, error)
	Identity(ctx context.Context, accessToken string) (*domain.User, error)
	UpdateMe(ctx context.Context, id uint, in UpdateMeInput) (*domain.User, error)
	Refresh(ctx context.Context, refreshToken string) (*UserLoginResult, error)
}
