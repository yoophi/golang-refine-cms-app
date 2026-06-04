package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/yoophi/refine-cms/backend/internal/bootstrap"
	"github.com/yoophi/refine-cms/backend/internal/config"
)

func main() {
	cfg := config.Load()

	injector := bootstrap.NewInjector(cfg)

	logger := do.MustInvoke[*zap.Logger](injector)
	defer func() { _ = logger.Sync() }()

	engine := do.MustInvoke[*gin.Engine](injector)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: engine,
	}

	// 서버 비동기 기동
	go func() {
		logger.Info("HTTP 서버 시작", zap.String("addr", srv.Addr), zap.String("db_driver", cfg.DB.Driver))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("서버 기동 실패", zap.Error(err))
		}
	}()

	// SIGINT/SIGTERM 수신 시 graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("종료 신호 수신, graceful shutdown 시작")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown 실패", zap.Error(err))
	}
	// DI 컨테이너에 등록된 종료 훅(예: DB) 실행
	if err := injector.Shutdown(); err != nil {
		logger.Error("컨테이너 종료 실패", zap.Error(err))
	}
	logger.Info("서버 정상 종료")
}
