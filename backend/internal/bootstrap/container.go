package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/yoophi/refine-cms/backend/internal/config"
	httpadapter "github.com/yoophi/refine-cms/backend/internal/adapter/handler/http"
	"github.com/yoophi/refine-cms/backend/internal/adapter/storage"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
	"github.com/yoophi/refine-cms/backend/internal/core/service"
)

// dbHandle 은 *sqlx.DB 와 dialect 드라이버명을 함께 담아 DI 로 전달하기 위한 래퍼이다.
type dbHandle struct {
	DB     *sqlx.DB
	Driver string
}

// NewInjector 는 samber/do 컨테이너를 구성한다.
// 의존성은 어댑터 -> 서비스 -> 포트 -> 도메인 순으로 안쪽을 향해 주입된다.
func NewInjector(cfg *config.Config) *do.Injector {
	i := do.New()

	do.ProvideValue(i, cfg)

	// --- 인프라: 로거 ---
	do.Provide(i, func(in *do.Injector) (*zap.Logger, error) {
		c := do.MustInvoke[*config.Config](in)
		if c.AppEnv == "production" {
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

	// --- 인바운드 어댑터: HTTP 핸들러 ---
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

	// --- 인바운드 어댑터: gin 엔진 ---
	do.Provide(i, func(in *do.Injector) (*gin.Engine, error) {
		if do.MustInvoke[*config.Config](in).AppEnv == "production" {
			gin.SetMode(gin.ReleaseMode)
		}
		return httpadapter.NewRouter(
			do.MustInvoke[*zap.Logger](in),
			httpadapter.Handlers{
				Category: do.MustInvoke[*httpadapter.CategoryHandler](in),
				Tag:      do.MustInvoke[*httpadapter.TagHandler](in),
				Post:     do.MustInvoke[*httpadapter.PostHandler](in),
				Comment:  do.MustInvoke[*httpadapter.CommentHandler](in),
			},
		), nil
	})

	return i
}
