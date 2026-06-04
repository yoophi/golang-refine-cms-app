package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CommentHandler 는 관리자 댓글 엔드포인트를 담당한다.
type CommentHandler struct {
	svc port.CommentService
}

func NewCommentHandler(svc port.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/comments")
	g.GET("", RequirePermission("comments", "list"), h.list)
	g.POST("", RequirePermission("comments", "create"), h.create)
	g.GET("/:id", RequirePermission("comments", "show"), h.get)
	g.PATCH("/:id", RequirePermission("comments", "edit"), h.update)
	g.DELETE("/:id", RequirePermission("comments", "delete"), h.delete)
}

func (h *CommentHandler) list(c *gin.Context) {
	q, err := newListQuery(c).eqInt("postId", "postId").eqStr("status", "status").like("authorName_like", "authorName").build()
	if err != nil {
		respondBadRequest(c, err)
		return
	}
	items, total, err := h.svc.Query(c.Request.Context(), q)
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]commentResponse, 0, len(items))
	for i := range items {
		res = append(res, newCommentResponse(&items[i]))
	}
	setTotalCount(c, total)
	c.JSON(http.StatusOK, res)
}

func (h *CommentHandler) get(c *gin.Context) {
	id, ok := parseID(c)
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

func (h *CommentHandler) create(c *gin.Context) {
	var req commentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	out, err := h.svc.Create(c.Request.Context(), port.CreateCommentInput{
		PostID:      req.PostID,
		ParentID:    req.ParentID,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
		Content:     req.Content,
		Status:      domain.CommentStatus(req.Status),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newCommentResponse(out))
}

func (h *CommentHandler) update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	p, err := newPatcher(c)
	if err != nil {
		respondBadRequest(c, err)
		return
	}
	cur, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	in := port.UpdateCommentInput{
		Content: cur.Content,
		Status:  cur.Status,
	}
	p.apply("content", &in.Content)
	p.apply("status", &in.Status)
	if p.err != nil {
		respondBadRequest(c, p.err)
		return
	}

	out, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newCommentResponse(out))
}

func (h *CommentHandler) delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
