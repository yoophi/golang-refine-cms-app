package security

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// BcryptHasher 는 port.PasswordHasher 의 bcrypt 구현이다.
type BcryptHasher struct {
	cost int
}

var _ port.PasswordHasher = (*BcryptHasher)(nil)

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: bcrypt.DefaultCost}
}

func (h *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", errors.Wrap(err, "비밀번호 해시")
	}
	return string(b), nil
}

// Compare 는 불일치 시 domain.ErrUnauthorized 를 반환한다(상세 사유는 노출하지 않음).
func (h *BcryptHasher) Compare(hash, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return domain.ErrUnauthorized
	}
	return nil
}
