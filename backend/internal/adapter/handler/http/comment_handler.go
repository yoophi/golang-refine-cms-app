package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CommentHandler 는 댓글 HTTP 엔드포인트를 담당한다.
type CommentHandler struct {
	svc port.CommentService
}

func NewCommentHandler(svc port.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/comments")
	g.POST("", h.create)
	g.GET("/:id", h.get)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)

	// 특정 게시글의 댓글 목록 (중첩 라우트)
	rg.GET("/posts/:id/comments", h.listByPost)
}

func (h *CommentHandler) create(c *gin.Context) {
	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	out, err := h.svc.Create(c.Request.Context(), port.CreateCommentInput{
		PostID:      req.PostID,
		ParentID:    req.ParentID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		Content:     req.Content,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newCommentResponse(out))
}

func (h *CommentHandler) get(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newCommentResponse(out))
}

func (h *CommentHandler) listByPost(c *gin.Context) {
	postID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.ListByPost(c.Request.Context(), postID)
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]commentResponse, 0, len(items))
	for i := range items {
		res = append(res, newCommentResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *CommentHandler) update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req updateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	out, err := h.svc.Update(c.Request.Context(), id, port.UpdateCommentInput{
		Content: req.Content,
		Status:  domain.CommentStatus(req.Status),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newCommentResponse(out))
}

func (h *CommentHandler) delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
