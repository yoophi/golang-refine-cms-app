package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// PostHandler 는 공개 게시글 조회(읽기 전용) 엔드포인트를 담당한다.
// 생성/수정/삭제 등 관리는 관리자 API(/admin/api/v1)에서만 제공한다.
type PostHandler struct {
	svc port.PostService
}

func NewPostHandler(svc port.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/posts")
	g.GET("", h.list)
	g.GET("/:id", h.get)
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
