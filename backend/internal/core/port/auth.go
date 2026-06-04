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

// AdminClaims 는 액세스 토큰에 담기는 관리자 신원 클레임이다.
type AdminClaims struct {
	UserID uint
	Email  string
	Name   string
	Role   domain.Role
}

// TokenManager 는 JWT 발급/검증 계약이다(인프라 어댑터가 구현).
type TokenManager interface {
	// Issue 는 사용자에 대한 액세스/리프레시 토큰을 발급한다.
	Issue(u *domain.AdminUser) (accessToken, refreshToken string, err error)
	// ParseAccess 는 액세스 토큰을 검증하고 클레임을 반환한다(무효/만료 시 domain.ErrUnauthorized).
	ParseAccess(token string) (*AdminClaims, error)
	// ParseRefresh 는 리프레시 토큰을 검증하고 사용자 ID 를 반환한다.
	ParseRefresh(token string) (uint, error)
}

// PasswordHasher 는 비밀번호 해시/검증 계약이다.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	// Compare 는 불일치 시 domain.ErrUnauthorized 를 반환한다.
	Compare(hash, plain string) error
}

// --- 입력 포트(Driving) ---

// LoginResult 는 로그인/리프레시 결과이다.
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
