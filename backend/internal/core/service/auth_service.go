package service

import (
	"context"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

type authService struct {
	repo   port.AdminUserRepository
	tokens port.TokenManager
	hasher port.PasswordHasher
}

// NewAuthService 는 AuthService 구현을 생성한다.
func NewAuthService(repo port.AdminUserRepository, tokens port.TokenManager, hasher port.PasswordHasher) port.AuthService {
	return &authService{repo: repo, tokens: tokens, hasher: hasher}
}

func (s *authService) Login(ctx context.Context, email, password string) (*port.LoginResult, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return nil, domain.ErrUnauthorized
	}

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// 존재하지 않는 계정도 자격 증명 불일치로 일원화(계정 존재 여부 노출 방지).
		return nil, domain.ErrUnauthorized
	}
	if err := s.hasher.Compare(u.PasswordHash, password); err != nil {
		return nil, domain.ErrUnauthorized
	}
	return s.issue(u)
}

func (s *authService) Identity(ctx context.Context, accessToken string) (*domain.AdminUser, error) {
	id, _, err := s.tokens.ParseAccess(accessToken)
	if err != nil {
		return nil, err
	}
	// 토큰이 유효해도 계정이 삭제/변경됐을 수 있으므로 최신 상태를 조회한다(role 도 DB 기준).
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return u, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*port.LoginResult, error) {
	id, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return nil, err
	}
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return s.issue(u)
}

func (s *authService) issue(u *domain.AdminUser) (*port.LoginResult, error) {
	access, refresh, err := s.tokens.Issue(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}
	return &port.LoginResult{AccessToken: access, RefreshToken: refresh, User: u}, nil
}
