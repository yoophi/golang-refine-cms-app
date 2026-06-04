package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// TagHandler 는 태그 HTTP 엔드포인트를 담당한다.
type TagHandler struct {
	svc port.TagService
}

func NewTagHandler(svc port.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

func (h *TagHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/tags")
	g.POST("", h.create)
	g.GET("", h.list)
	g.GET("/:id", h.get)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *TagHandler) create(c *gin.Context) {
	var req createTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	out, err := h.svc.Create(c.Request.Context(), port.CreateTagInput{Name: req.Name, Slug: req.Slug})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newTagResponse(out))
}

func (h *TagHandler) list(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]tagResponse, 0, len(items))
	for i := range items {
		res = append(res, newTagResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *TagHandler) get(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newTagResponse(out))
}

func (h *TagHandler) update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req updateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	out, err := h.svc.Update(c.Request.Context(), id, port.UpdateTagInput{Name: req.Name, Slug: req.Slug})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newTagResponse(out))
}

func (h *TagHandler) delete(c *gin.Context) {
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
