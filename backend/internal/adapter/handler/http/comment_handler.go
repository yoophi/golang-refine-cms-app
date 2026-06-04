package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CommentHandler 는 공개 댓글 HTTP 엔드포인트를 담당한다.
type CommentHandler struct {
	svc port.CommentService
}

func NewCommentHandler(svc port.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// register 는 댓글 라우트를 등록한다. auth 미들웨어는 작성/수정/삭제(로그인 회원 + 소유권)에만 적용한다.
func (h *CommentHandler) register(rg *gin.RouterGroup, auth gin.HandlerFunc) {
	// 목록은 공개(응답에 user_id 포함)
	rg.GET("/posts/:id/comments", h.listByPost)

	g := rg.Group("/comments")
	g.GET("/:id", h.get) // 단건 조회는 공개
	g.POST("", auth, h.create)
	g.PATCH("/:id", auth, h.update)
	g.DELETE("/:id", auth, h.delete)
}

func (h *CommentHandler) create(c *gin.Context) {
	u, ok := userIdentityOf(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	// 작성자는 토큰 회원에서 유도한다(클라이언트 입력 무시).
	uid := u.ID
	out, err := h.svc.Create(c.Request.Context(), port.CreateCommentInput{
		PostID:      req.PostID,
		ParentID:    req.ParentID,
		UserID:      &uid,
		AuthorName:  u.Name,
		AuthorEmail: u.Email,
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
	u, ok := userIdentityOf(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req updateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	out, err := h.svc.UpdateOwnedContent(c.Request.Context(), id, u.ID, req.Content)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newCommentResponse(out))
}

func (h *CommentHandler) delete(c *gin.Context) {
	u, ok := userIdentityOf(c)
	if !ok {
		respondError(c, domain.ErrUnauthorized)
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteOwned(c.Request.Context(), id, u.ID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
