package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// TagHandler 는 관리자 태그 엔드포인트를 담당한다.
type TagHandler struct {
	svc port.TagService
}

func NewTagHandler(svc port.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

func (h *TagHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/tags")
	g.GET("", RequirePermission("tags", "list"), h.list)
	g.POST("", RequirePermission("tags", "create"), h.create)
	g.GET("/:id", RequirePermission("tags", "show"), h.get)
	g.PATCH("/:id", RequirePermission("tags", "edit"), h.update)
	g.DELETE("/:id", RequirePermission("tags", "delete"), h.delete)
}

func (h *TagHandler) list(c *gin.Context) {
	q, err := newListQuery(c).like("name_like", "name").build()
	if err != nil {
		respondBadRequest(c, err)
		return
	}
	items, total, err := h.svc.Query(c.Request.Context(), q)
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]tagResponse, 0, len(items))
	for i := range items {
		res = append(res, newTagResponse(&items[i]))
	}
	setTotalCount(c, total)
	c.JSON(http.StatusOK, res)
}

func (h *TagHandler) get(c *gin.Context) {
	id, ok := parseID(c)
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

func (h *TagHandler) create(c *gin.Context) {
	var req tagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	out, err := h.svc.Create(c.Request.Context(), port.CreateTagInput{Name: req.Name, Slug: req.Slug})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newTagResponse(out))
}

func (h *TagHandler) update(c *gin.Context) {
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

	in := port.UpdateTagInput{Name: cur.Name, Slug: cur.Slug}
	p.apply("name", &in.Name)
	p.apply("slug", &in.Slug)
	if p.err != nil {
		respondBadRequest(c, p.err)
		return
	}

	out, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newTagResponse(out))
}

func (h *TagHandler) delete(c *gin.Context) {
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
