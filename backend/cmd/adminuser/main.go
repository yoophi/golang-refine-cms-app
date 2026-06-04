// Command adminuser 는 관리자 계정(admin_users)을 생성/갱신하는 CLI 도구이다.
//
// 사용 예:
//
//	go run ./cmd/adminuser -email test@example.com -password secret11 -role superadmin
//
// DB 연결은 애플리케이션과 동일하게 환경변수(.env)에서 로딩한다(SQLITE_PATH/DB_DRIVER 등).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yoophi/refine-cms/backend/internal/adapter/security"
	"github.com/yoophi/refine-cms/backend/internal/adapter/storage"
	"github.com/yoophi/refine-cms/backend/internal/config"
	"github.com/yoophi/refine-cms/backend/internal/core/domain"
)

func main() {
	email := flag.String("email", "", "관리자 이메일(필수)")
	password := flag.String("password", "", "비밀번호(필수)")
	name := flag.String("name", "", "표시 이름(생략 시 이메일)")
	role := flag.String("role", "superadmin", "역할: superadmin | editor | viewer")
	flag.Parse()

	if *email == "" || *password == "" {
		fatalf("email 과 password 는 필수입니다.\n사용: go run ./cmd/adminuser -email <email> -password <pw> [-role superadmin|editor|viewer] [-name <name>]")
	}
	r := domain.Role(*role)
	if !r.Valid() {
		fatalf("유효하지 않은 role: %q (superadmin|editor|viewer 중 하나)", *role)
	}

	cfg := config.Load()
	db, driver, err := storage.NewDB(cfg.DB)
	if err != nil {
		fatalf("DB 연결 실패: %v", err)
	}
	if err := storage.Migrate(db, driver); err != nil {
		fatalf("마이그레이션 실패: %v", err)
	}

	repo := storage.NewAdminUserRepository(db, driver)
	hasher := security.NewBcryptHasher()
	ctx := context.Background()

	emailNorm := strings.ToLower(strings.TrimSpace(*email))
	if _, err := repo.GetByEmail(ctx, emailNorm); err == nil {
		fatalf("이미 존재하는 이메일입니다: %s", emailNorm)
	}

	hash, err := hasher.Hash(*password)
	if err != nil {
		fatalf("비밀번호 해시 실패: %v", err)
	}
	displayName := strings.TrimSpace(*name)
	if displayName == "" {
		displayName = emailNorm
	}

	u := &domain.AdminUser{
		Name:         displayName,
		Email:        emailNorm,
		Role:         r,
		PasswordHash: hash,
	}
	if err := repo.Create(ctx, u); err != nil {
		fatalf("관리자 생성 실패: %v", err)
	}

	fmt.Printf("관리자 계정 등록됨: id=%d email=%s role=%s (driver=%s)\n", u.ID, u.Email, u.Role, driver)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
