package http

import (
	"time"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// ---- 회원 요청 DTO ----

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type userLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// updateMeRequest: 모든 필드 선택(부분 수정). 포인터로 '제공 여부'를 구분한다.
type updateMeRequest struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Password        *string `json:"password"`
	CurrentPassword *string `json:"currentPassword"`
}

type userRefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ---- 회원 응답 DTO ----

// userResponse: 회원 모델(비밀번호 해시 제외). 필드는 공개 API 규약(snake_case).
type userResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// 인증 응답 봉투(토큰 필드는 계약상 camelCase).
type userAuthResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken,omitempty"`
	User         userResponse `json:"user"`
}

type userRefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
