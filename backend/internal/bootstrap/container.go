package bootstrap

import (
	"context"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"go.uber.org/zap"

	adminhttp "github.com/yoophi/refine-cms/backend/internal/adapter/handler/admin"
	httpadapter "github.com/yoophi/refine-cms/backend/internal/adapter/handler/http"
	"github.com/yoophi/refine-cms/backend/internal/adapter/security"
	"github.com/yoophi/refine-cms/backend/internal/adapter/storage"
	"github.com/yoophi/refine-cms/backend/internal/config"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
	"github.com/yoophi/refine-cms/backend/internal/core/service"
)

// dbHandle 은 *sqlx.DB 와 dialect 드라이버명을 함께 담아 DI 로 전달하기 위한 래퍼이다.
type dbHandle struct {
	DB     *sqlx.DB
	Driver string
}

// newJWTManager 는 동일 시크릿/TTL 로 audience 만 다른 토큰 매니저를 만든다(admin/user 분리).
func newJWTManager(c *config.Config, audience string) *security.JWTManager {
	return security.NewJWTManager(c.Admin.JWTSecret, c.Admin.AccessTTL, c.Admin.RefreshTTL, audience)
}

// NewInjector 는 samber/do 컨테이너를 구성한다.
// 의존성은 어댑터 -> 서비스 -> 포트 -> 도메인 순으로 안쪽을 향해 주입된다.
func NewInjector(cfg *config.Config) *do.Injector {
	i := do.New()

	do.ProvideValue(i, cfg)

	// --- 인프라: 로거 ---
	do.Provide(i, func(in *do.Injector) (*zap.Logger, error) {
		if do.MustInvoke[*config.Config](in).AppEnv == "production" {
			return zap.NewProduction()
		}
		return zap.NewDevelopment()
	})

	// --- 인프라: DB 연결 + 마이그레이션 ---
	do.Provide(i, func(in *do.Injector) (*dbHandle, error) {
		c := do.MustInvoke[*config.Config](in)
		db, driver, err := storage.NewDB(c.DB)
		if err != nil {
			return nil, err
		}
		if err := storage.Migrate(db, driver); err != nil {
			return nil, err
		}
		return &dbHandle{DB: db, Driver: driver}, nil
	})

	// --- 인프라: 비밀번호 해시 ---
	// TokenManager 는 audience(admin/user)별로 분리해야 하므로 공용 provider 로 두지 않고
	// 각 인증 서비스 provider 에서 직접 생성한다(동일 시크릿, 다른 audience).
	do.Provide(i, func(in *do.Injector) (port.PasswordHasher, error) {
		return security.NewBcryptHasher(), nil
	})

	// --- 아웃바운드 어댑터: 리포지토리(출력 포트 구현) ---
	do.Provide(i, func(in *do.Injector) (port.CategoryRepository, error) {
		h := do.MustInvoke[*dbHandle](in)
		return storage.NewCategoryRepository(h.DB, h.Driver), nil
	})
	do.Provide(i, func(in *do.Injector) (port.TagRepository, error) {
		h := do.MustInvoke[*dbHandle](in)
		return storage.NewTagRepository(h.DB, h.Driver), nil
	})
	do.Provide(i, func(in *do.Injector) (port.PostRepository, error) {
		h := do.MustInvoke[*dbHandle](in)
		return storage.NewPostRepository(h.DB, h.Driver), nil
	})
	do.Provide(i, func(in *do.Injector) (port.CommentRepository, error) {
		h := do.MustInvoke[*dbHandle](in)
		return storage.NewCommentRepository(h.DB, h.Driver), nil
	})
	// 관리자 사용자 리포지토리(구체 타입으로도 제공: 시드에 Count/Create 필요).
	do.Provide(i, func(in *do.Injector) (*storage.AdminUserRepository, error) {
		h := do.MustInvoke[*dbHandle](in)
		return storage.NewAdminUserRepository(h.DB, h.Driver), nil
	})
	do.Provide(i, func(in *do.Injector) (port.AdminUserRepository, error) {
		return do.MustInvoke[*storage.AdminUserRepository](in), nil
	})
	do.Provide(i, func(in *do.Injector) (port.UserRepository, error) {
		h := do.MustInvoke[*dbHandle](in)
		return storage.NewUserRepository(h.DB, h.Driver), nil
	})

	// --- 코어: 서비스(입력 포트 구현) ---
	do.Provide(i, func(in *do.Injector) (port.CategoryService, error) {
		return service.NewCategoryService(do.MustInvoke[port.CategoryRepository](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (port.TagService, error) {
		return service.NewTagService(do.MustInvoke[port.TagRepository](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (port.PostService, error) {
		return service.NewPostService(do.MustInvoke[port.PostRepository](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (port.CommentService, error) {
		return service.NewCommentService(
			do.MustInvoke[port.CommentRepository](in),
			do.MustInvoke[port.PostRepository](in),
		), nil
	})
	do.Provide(i, func(in *do.Injector) (port.AuthService, error) {
		tokens := newJWTManager(do.MustInvoke[*config.Config](in), security.AudienceAdmin)
		return service.NewAuthService(
			do.MustInvoke[port.AdminUserRepository](in),
			tokens,
			do.MustInvoke[port.PasswordHasher](in),
		), nil
	})
	do.Provide(i, func(in *do.Injector) (port.UserAuthService, error) {
		tokens := newJWTManager(do.MustInvoke[*config.Config](in), security.AudienceUser)
		return service.NewUserAuthService(
			do.MustInvoke[port.UserRepository](in),
			tokens,
			do.MustInvoke[port.PasswordHasher](in),
		), nil
	})

	// --- 인바운드 어댑터: 사용자(public) HTTP 핸들러 ---
	do.Provide(i, func(in *do.Injector) (*httpadapter.CategoryHandler, error) {
		return httpadapter.NewCategoryHandler(do.MustInvoke[port.CategoryService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*httpadapter.TagHandler, error) {
		return httpadapter.NewTagHandler(do.MustInvoke[port.TagService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*httpadapter.PostHandler, error) {
		return httpadapter.NewPostHandler(do.MustInvoke[port.PostService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*httpadapter.CommentHandler, error) {
		return httpadapter.NewCommentHandler(do.MustInvoke[port.CommentService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*httpadapter.UserAuthHandler, error) {
		return httpadapter.NewUserAuthHandler(do.MustInvoke[port.UserAuthService](in)), nil
	})

	// --- 인바운드 어댑터: 관리자(admin) HTTP 핸들러 ---
	do.Provide(i, func(in *do.Injector) (*adminhttp.AuthHandler, error) {
		return adminhttp.NewAuthHandler(do.MustInvoke[port.AuthService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*adminhttp.PostHandler, error) {
		return adminhttp.NewPostHandler(do.MustInvoke[port.PostService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*adminhttp.CategoryHandler, error) {
		return adminhttp.NewCategoryHandler(do.MustInvoke[port.CategoryService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*adminhttp.TagHandler, error) {
		return adminhttp.NewTagHandler(do.MustInvoke[port.TagService](in)), nil
	})
	do.Provide(i, func(in *do.Injector) (*adminhttp.CommentHandler, error) {
		return adminhttp.NewCommentHandler(do.MustInvoke[port.CommentService](in)), nil
	})

	// --- 인바운드 어댑터: gin 엔진(전역 미들웨어 + 라우트 그룹 조립) ---
	do.Provide(i, func(in *do.Injector) (*gin.Engine, error) {
		c := do.MustInvoke[*config.Config](in)
		logger := do.MustInvoke[*zap.Logger](in)

		// 개발용 기본 관리자 시드(admin_users 가 비어있을 때만). 운영에서는 알려진 기본 계정 생성을
		// 막기 위해 자동 시드를 비활성화한다(계정은 cmd/adminuser 로 명시 생성).
		if c.AppEnv == "production" {
			logger.Info("운영 환경: 기본 관리자 자동 시드 비활성화 (cmd/adminuser 로 계정 생성)")
		} else {
			seeded, err := storage.SeedDefaultAdmins(
				context.Background(),
				do.MustInvoke[*storage.AdminUserRepository](in),
				do.MustInvoke[port.PasswordHasher](in),
				c.Admin.SeedPassword,
			)
			if err != nil {
				return nil, err
			}
			if seeded {
				logger.Warn("기본 관리자 계정 시드됨(개발용) — 운영 전 비밀번호/계정 교체 필요",
					zap.String("accounts", "admin@example.com / editor@example.com / viewer@example.com"))
			}
		}

		if c.AppEnv == "production" {
			gin.SetMode(gin.ReleaseMode)
		}

		r := gin.New()
		// 전역 미들웨어: 로깅 / panic 복구 / CORS(프리플라이트 포함).
		r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		r.Use(ginzap.RecoveryWithZap(logger, true))
		r.Use(adminhttp.CORS(c.Admin.CORSOrigins))

		r.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// 사용자(public) API: /api/v1 (회원 인증 + 댓글 소유권)
		httpadapter.RegisterRoutes(r.Group("/api/v1"), httpadapter.Handlers{
			Category: do.MustInvoke[*httpadapter.CategoryHandler](in),
			Tag:      do.MustInvoke[*httpadapter.TagHandler](in),
			Post:     do.MustInvoke[*httpadapter.PostHandler](in),
			Comment:  do.MustInvoke[*httpadapter.CommentHandler](in),
			UserAuth: do.MustInvoke[*httpadapter.UserAuthHandler](in),
		}, do.MustInvoke[port.UserAuthService](in))

		// 관리자 API: /admin/api/v1 (인증 + ACL)
		adminhttp.RegisterRoutes(
			r.Group("/admin/api/v1"),
			do.MustInvoke[port.AuthService](in),
			adminhttp.Handlers{
				Auth:     do.MustInvoke[*adminhttp.AuthHandler](in),
				Post:     do.MustInvoke[*adminhttp.PostHandler](in),
				Category: do.MustInvoke[*adminhttp.CategoryHandler](in),
				Tag:      do.MustInvoke[*adminhttp.TagHandler](in),
				Comment:  do.MustInvoke[*adminhttp.CommentHandler](in),
			},
		)

		return r, nil
	})

	return i
}
