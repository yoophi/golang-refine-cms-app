package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CategoryHandler 는 관리자 카테고리 엔드포인트를 담당한다.
type CategoryHandler struct {
	svc port.CategoryService
}

func NewCategoryHandler(svc port.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/categories")
	g.GET("", RequirePermission("categories", "list"), h.list)
	g.POST("", RequirePermission("categories", "create"), h.create)
	g.GET("/:id", RequirePermission("categories", "show"), h.get)
	g.PATCH("/:id", RequirePermission("categories", "edit"), h.update)
	g.DELETE("/:id", RequirePermission("categories", "delete"), h.delete)
}

func (h *CategoryHandler) list(c *gin.Context) {
	q, err := newListQuery(c).like("name_like", "name").eqInt("parentId", "parentId").build()
	if err != nil {
		respondBadRequest(c, err)
		return
	}
	items, total, err := h.svc.Query(c.Request.Context(), q)
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]categoryResponse, 0, len(items))
	for i := range items {
		res = append(res, newCategoryResponse(&items[i]))
	}
	setTotalCount(c, total)
	c.JSON(http.StatusOK, res)
}

func (h *CategoryHandler) get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	out, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newCategoryResponse(out))
}

func (h *CategoryHandler) create(c *gin.Context) {
	var req categoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err)
		return
	}
	out, err := h.svc.Create(c.Request.Context(), port.CreateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, newCategoryResponse(out))
}

func (h *CategoryHandler) update(c *gin.Context) {
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

	in := port.UpdateCategoryInput{
		Name:        cur.Name,
		Slug:        cur.Slug,
		Description: cur.Description,
		ParentID:    cur.ParentID,
	}
	p.apply("name", &in.Name)
	p.apply("slug", &in.Slug)
	p.apply("description", &in.Description)
	p.apply("parentId", &in.ParentID)
	if p.err != nil {
		respondBadRequest(c, p.err)
		return
	}

	out, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newCategoryResponse(out))
}

func (h *CategoryHandler) delete(c *gin.Context) {
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
