package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// AuthHandler 는 관리자 인증 엔드포인트(/auth/*)를 담당한다.
type AuthHandler struct {
	svc port.AuthService
}

func NewAuthHandler(svc port.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type loginResponse struct {
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken,omitempty"`
	User         adminUserResponse `json:"user"`
}

type refreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// login (공개): 이메일/비밀번호 → 토큰 + 신원.
func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// 자격 증명 실패는 사유를 구분하지 않고 401 + 명시 메시지로 응답한다.
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "이메일 또는 비밀번호가 올바르지 않습니다."})
		return
	}
	c.JSON(http.StatusOK, loginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User:         newAdminUserResponse(res.User),
	})
}

// me (보호): 현재 토큰의 관리자 신원.
func (h *AuthHandler) me(c *gin.Context) {
	u, ok := identityOf(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	c.JSON(http.StatusOK, newAdminUserResponse(u))
}

// refresh (공개): 리프레시 토큰으로 액세스 토큰 재발급.
func (h *AuthHandler) refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, refreshResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// logout (보호): 스테이트리스 JWT 이므로 서버 상태 변경 없이 200 을 반환한다.
// (대시보드는 응답과 무관하게 로컬 토큰을 폐기한다. 서버측 무효화가 필요하면 블랙리스트 도입.)
func (h *AuthHandler) logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
