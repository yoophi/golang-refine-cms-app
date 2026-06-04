package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

// identityKey 는 인증된 관리자 신원을 gin 컨텍스트에 저장하는 키이다.
const identityKey = "adminIdentity"

func setIdentity(c *gin.Context, u *domain.AdminUser) {
	c.Set(identityKey, u)
}

// identityOf 는 컨텍스트에 저장된 인증 신원을 반환한다(Authenticate 미들웨어가 설정).
func identityOf(c *gin.Context) (*domain.AdminUser, bool) {
	v, ok := c.Get(identityKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*domain.AdminUser)
	return u, ok
}

// adminUserResponse 는 AdminUser 의 camelCase 표현이다(비밀번호 해시는 노출하지 않음).
type adminUserResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	Avatar      string   `json:"avatar,omitempty"`
}

func newAdminUserResponse(u *domain.AdminUser) adminUserResponse {
	return adminUserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        string(u.Role),
		Permissions: u.Permissions(),
		Avatar:      u.Avatar,
	}
}
