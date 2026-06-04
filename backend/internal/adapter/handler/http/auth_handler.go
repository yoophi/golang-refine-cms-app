package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// UserAuthHandler 는 공개 회원 인증 엔드포인트(/auth/*)를 담당한다.
type UserAuthHandler struct {
	svc port.UserAuthService
}

func NewUserAuthHandler(svc port.UserAuthService) *UserAuthHandler {
	return &UserAuthHandler{svc: svc}
}

// register (공개): 회원 가입 → 토큰 + 회원.
func (h *UserAuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	res, err := h.svc.Register(c.Request.Context(), port.RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		respondError(c, err) // 이메일 중복 → 409
		return
	}
	c.JSON(http.StatusCreated, userAuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User:         newUserResponse(res.User),
	})
}

// login (공개): 로그인.
func (h *UserAuthHandler) login(c *gin.Context) {
	var req userLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	res, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "이메일 또는 비밀번호가 올바르지 않습니다."})
		return
	}
	c.JSON(http.StatusOK, userAuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		User:         newUserResponse(res.User),
	})
}

// me (보호): 현재 회원 정보.
func (h *UserAuthHandler) me(c *gin.Context) {
	u, ok := userIdentityOf(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	c.JSON(http.StatusOK, newUserResponse(u))
}

// updateMe (보호): 회원 정보 부분 수정.
func (h *UserAuthHandler) updateMe(c *gin.Context) {
	u, ok := userIdentityOf(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	out, err := h.svc.UpdateMe(c.Request.Context(), u.ID, port.UpdateMeInput{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		CurrentPassword: req.CurrentPassword,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newUserResponse(out))
}

// refresh (공개): 리프레시 토큰으로 액세스 토큰 재발급.
func (h *UserAuthHandler) refresh(c *gin.Context) {
	var req userRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	res, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, userRefreshResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// logout (보호): 스테이트리스 JWT 이므로 200 만 반환.
func (h *UserAuthHandler) logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
