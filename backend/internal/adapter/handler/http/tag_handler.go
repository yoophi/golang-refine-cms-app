package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// TagHandler 는 공개 태그 조회(읽기 전용) 엔드포인트를 담당한다.
// 생성/수정/삭제 등 관리는 관리자 API(/admin/api/v1)에서만 제공한다.
type TagHandler struct {
	svc port.TagService
}

func NewTagHandler(svc port.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

func (h *TagHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/tags")
	g.GET("", h.list)
	g.GET("/:id", h.get)
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
