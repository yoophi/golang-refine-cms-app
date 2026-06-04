package security

import (
	"slices"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// 토큰 audience: 회원 토큰과 관리자 토큰을 분리해 상호 접근을 차단한다(동일 시크릿).
const (
	AudienceAdmin = "admin"
	AudienceUser  = "user"
)

const (
	typeAccess  = "access"
	typeRefresh = "refresh"
)

// JWTManager 는 port.TokenManager 의 golang-jwt/v4(HS256) 구현이다.
// audience 로 토큰 용도(admin/user)를 구분하고, typ 로 access/refresh 를 구분한다.
type JWTManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	audience   string
}

var _ port.TokenManager = (*JWTManager)(nil)

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration, audience string) *JWTManager {
	return &JWTManager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL, audience: audience}
}

type claims struct {
	Role string `json:"role,omitempty"`
	Typ  string `json:"typ"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Issue(subject uint, role string) (string, string, error) {
	now := time.Now()
	sub := strconv.FormatUint(uint64(subject), 10)
	aud := jwt.ClaimStrings{m.audience}

	access, err := m.sign(claims{
		Role: role,
		Typ:  typeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Audience:  aud,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	})
	if err != nil {
		return "", "", errors.Wrap(err, "액세스 토큰 서명")
	}
	refresh, err := m.sign(claims{
		Typ: typeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Audience:  aud,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
		},
	})
	if err != nil {
		return "", "", errors.Wrap(err, "리프레시 토큰 서명")
	}
	return access, refresh, nil
}

func (m *JWTManager) ParseAccess(token string) (uint, string, error) {
	c, err := m.parse(token, typeAccess)
	if err != nil {
		return 0, "", err
	}
	id, err := subjectID(c)
	if err != nil {
		return 0, "", err
	}
	return id, c.Role, nil
}

func (m *JWTManager) ParseRefresh(token string) (uint, error) {
	c, err := m.parse(token, typeRefresh)
	if err != nil {
		return 0, err
	}
	return subjectID(c)
}

func (m *JWTManager) sign(c claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)
}

// parse 는 서명/만료 + typ + audience 를 검증한다. 어떤 실패든 domain.ErrUnauthorized 로 일원화.
func (m *JWTManager) parse(token, wantTyp string) (*claims, error) {
	var c claims
	t, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return m.secret, nil
	})
	if err != nil || !t.Valid {
		return nil, domain.ErrUnauthorized
	}
	if c.Typ != wantTyp {
		return nil, domain.ErrUnauthorized // access/refresh 혼용 차단
	}
	if !slices.Contains(c.Audience, m.audience) {
		return nil, domain.ErrUnauthorized // user/admin 토큰 상호 사용 차단
	}
	return &c, nil
}

func subjectID(c *claims) (uint, error) {
	id, err := strconv.ParseUint(c.Subject, 10, 64)
	if err != nil {
		return 0, domain.ErrUnauthorized
	}
	return uint(id), nil
}
