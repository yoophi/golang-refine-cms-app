package security

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// JWTManager 는 port.TokenManager 의 golang-jwt/v4(HS256) 구현이다.
type JWTManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

var _ port.TokenManager = (*JWTManager)(nil)

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

// accessClaims 는 액세스 토큰 페이로드이다. sub=관리자 id.
type accessClaims struct {
	Role  string `json:"role"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

// refreshClaims 는 리프레시 토큰 페이로드이다(typ 으로 액세스와 구분).
type refreshClaims struct {
	Typ string `json:"typ"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Issue(u *domain.AdminUser) (string, string, error) {
	now := time.Now()
	sub := strconv.FormatUint(uint64(u.ID), 10)

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims{
		Role:  string(u.Role),
		Email: u.Email,
		Name:  u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	})
	accessStr, err := access.SignedString(m.secret)
	if err != nil {
		return "", "", errors.Wrap(err, "액세스 토큰 서명")
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims{
		Typ: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
		},
	})
	refreshStr, err := refresh.SignedString(m.secret)
	if err != nil {
		return "", "", errors.Wrap(err, "리프레시 토큰 서명")
	}
	return accessStr, refreshStr, nil
}

func (m *JWTManager) ParseAccess(token string) (*port.AdminClaims, error) {
	var claims accessClaims
	if _, err := m.parse(token, &claims); err != nil {
		return nil, err
	}
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return &port.AdminClaims{
		UserID: uint(id),
		Email:  claims.Email,
		Name:   claims.Name,
		Role:   domain.Role(claims.Role),
	}, nil
}

func (m *JWTManager) ParseRefresh(token string) (uint, error) {
	var claims refreshClaims
	if _, err := m.parse(token, &claims); err != nil {
		return 0, err
	}
	if claims.Typ != "refresh" {
		return 0, domain.ErrUnauthorized
	}
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, domain.ErrUnauthorized
	}
	return uint(id), nil
}

// parse 는 서명/만료를 검증한다. 어떤 실패든 domain.ErrUnauthorized 로 일원화한다.
func (m *JWTManager) parse(token string, claims jwt.Claims) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return m.secret, nil
	})
	if err != nil || !t.Valid {
		return nil, domain.ErrUnauthorized
	}
	return t, nil
}
