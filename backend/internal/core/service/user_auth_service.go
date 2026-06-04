package service

import (
	"context"
	"strings"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

type userAuthService struct {
	repo   port.UserRepository
	tokens port.TokenManager
	hasher port.PasswordHasher
}

// NewUserAuthService 는 회원 인증 유스케이스 구현을 생성한다.
// tokens 는 회원 audience 로 구성된 TokenManager 여야 한다(관리자와 분리).
func NewUserAuthService(repo port.UserRepository, tokens port.TokenManager, hasher port.PasswordHasher) port.UserAuthService {
	return &userAuthService{repo: repo, tokens: tokens, hasher: hasher}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *userAuthService) Register(ctx context.Context, in port.RegisterUserInput) (*port.UserLoginResult, error) {
	email := normalizeEmail(in.Email)
	if email == "" || in.Password == "" || strings.TrimSpace(in.Name) == "" {
		return nil, domain.ErrInvalidInput
	}
	// 이메일 유일성: 이미 있으면 충돌.
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, domain.ErrConflict
	}

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	u := &domain.User{Email: email, Name: strings.TrimSpace(in.Name), PasswordHash: hash}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return s.issue(u)
}

func (s *userAuthService) Login(ctx context.Context, email, password string) (*port.UserLoginResult, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return nil, domain.ErrUnauthorized
	}
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrUnauthorized // 계정 존재 여부 노출 방지
	}
	if err := s.hasher.Compare(u.PasswordHash, password); err != nil {
		return nil, domain.ErrUnauthorized
	}
	return s.issue(u)
}

func (s *userAuthService) Identity(ctx context.Context, accessToken string) (*domain.User, error) {
	id, _, err := s.tokens.ParseAccess(accessToken)
	if err != nil {
		return nil, err
	}
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return u, nil
}

func (s *userAuthService) UpdateMe(ctx context.Context, id uint, in port.UpdateMeInput) (*domain.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return nil, domain.ErrInvalidInput
		}
		u.Name = name
	}
	if in.Email != nil {
		email := normalizeEmail(*in.Email)
		if email == "" {
			return nil, domain.ErrInvalidInput
		}
		if email != u.Email {
			if _, err := s.repo.GetByEmail(ctx, email); err == nil {
				return nil, domain.ErrConflict
			}
			u.Email = email
		}
	}
	if in.Password != nil {
		if *in.Password == "" {
			return nil, domain.ErrInvalidInput
		}
		// 비밀번호 변경 시 현재 비밀번호 검증(제공된 경우).
		if in.CurrentPassword == nil || s.hasher.Compare(u.PasswordHash, *in.CurrentPassword) != nil {
			return nil, domain.ErrUnauthorized
		}
		hash, err := s.hasher.Hash(*in.Password)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = hash
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userAuthService) Refresh(ctx context.Context, refreshToken string) (*port.UserLoginResult, error) {
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

func (s *userAuthService) issue(u *domain.User) (*port.UserLoginResult, error) {
	access, refresh, err := s.tokens.Issue(u.ID, "")
	if err != nil {
		return nil, err
	}
	return &port.UserLoginResult{AccessToken: access, RefreshToken: refresh, User: u}, nil
}
