package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// CategoryHandler 는 공개 카테고리 조회(읽기 전용) 엔드포인트를 담당한다.
// 생성/수정/삭제 등 관리는 관리자 API(/admin/api/v1)에서만 제공한다.
type CategoryHandler struct {
	svc port.CategoryService
}

func NewCategoryHandler(svc port.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/categories")
	g.GET("", h.list)
	g.GET("/:id", h.get)
}

func (h *CategoryHandler) list(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]categoryResponse, 0, len(items))
	for i := range items {
		res = append(res, newCategoryResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *CategoryHandler) get(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
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
