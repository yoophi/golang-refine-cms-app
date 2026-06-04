package http

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handlers 는 라우터 구성에 필요한 모든 핸들러를 모은다.
type Handlers struct {
	Category *CategoryHandler
	Tag      *TagHandler
	Post     *PostHandler
	Comment  *CommentHandler
}

// NewRouter 는 미들웨어와 모든 라우트가 구성된 gin 엔진을 생성한다.
func NewRouter(logger *zap.Logger, h Handlers) *gin.Engine {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	h.Category.register(api)
	h.Tag.register(api)
	h.Post.register(api)
	h.Comment.register(api)

	return r
}
