package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// PostHandler 는 게시글 HTTP 엔드포인트를 담당한다.
type PostHandler struct {
	svc port.PostService
}

func NewPostHandler(svc port.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/posts")
	g.POST("", h.create)
	g.GET("", h.list)
	g.GET("/:id", h.get)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *PostHandler) create(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	out, err := h.svc.Create(c.Request.Context(), port.CreatePostInput{
		Title:      req.Title,
		Slug:       req.Slug,
		Excerpt:    req.Excerpt,
		Content:    req.Content,
		Status:     domain.PostStatus(req.Status),
		CategoryID: req.CategoryID,
		TagIDs:     req.TagIDs,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newPostResponse(out))
}

func (h *PostHandler) list(c *gin.Context) {
	var f port.PostFilter
	if s := c.Query("status"); s != "" {
		st := domain.PostStatus(s)
		f.Status = &st
	}
	if cid := c.Query("category_id"); cid != "" {
		if v, err := strconv.ParseUint(cid, 10, 64); err == nil {
			id := uint(v)
			f.CategoryID = &id
		}
	}

	items, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]postResponse, 0, len(items))
	for i := range items {
		res = append(res, newPostResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *PostHandler) get(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newPostResponse(out))
}

func (h *PostHandler) update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	out, err := h.svc.Update(c.Request.Context(), id, port.UpdatePostInput{
		Title:      req.Title,
		Slug:       req.Slug,
		Excerpt:    req.Excerpt,
		Content:    req.Content,
		Status:     domain.PostStatus(req.Status),
		CategoryID: req.CategoryID,
		TagIDs:     req.TagIDs,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newPostResponse(out))
}

func (h *PostHandler) delete(c *gin.Context) {
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
