package storage

import (
	"context"

	"github.com/pkg/errors"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// SeedDefaultAdmins 는 admin_users 가 비어있을 때만 개발용 기본 계정을 생성한다.
// 운영에서는 별도 프로비저닝을 사용하고, 이 시드 비밀번호를 그대로 두지 말 것.
// 반환값: 시드를 수행했는지 여부.
func SeedDefaultAdmins(ctx context.Context, repo *AdminUserRepository, hasher port.PasswordHasher, password string) (bool, error) {
	n, err := repo.Count(ctx)
	if err != nil {
		return false, err
	}
	if n > 0 {
		return false, nil
	}

	hash, err := hasher.Hash(password)
	if err != nil {
		return false, err
	}

	defaults := []domain.AdminUser{
		{Name: "관리자", Email: "admin@example.com", Role: domain.RoleSuperadmin, PasswordHash: hash},
		{Name: "에디터", Email: "editor@example.com", Role: domain.RoleEditor, PasswordHash: hash},
		{Name: "뷰어", Email: "viewer@example.com", Role: domain.RoleViewer, PasswordHash: hash},
	}
	for i := range defaults {
		u := defaults[i]
		if err := repo.Create(ctx, &u); err != nil {
			return false, errors.Wrap(err, "기본 관리자 시드")
		}
	}
	return true, nil
}
