package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// PostHandler 는 관리자 게시글 엔드포인트를 담당한다.
type PostHandler struct {
	svc port.PostService
}

func NewPostHandler(svc port.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) register(rg *gin.RouterGroup) {
	g := rg.Group("/posts")
	g.GET("", RequirePermission("posts", "list"), h.list)
	g.POST("", RequirePermission("posts", "create"), h.create)
	g.GET("/:id", RequirePermission("posts", "show"), h.get)
	g.PATCH("/:id", RequirePermission("posts", "edit"), h.update)
	g.DELETE("/:id", RequirePermission("posts", "delete"), h.delete)
}

func (h *PostHandler) list(c *gin.Context) {
	q, err := newListQuery(c).like("title_like", "title").eqStr("status", "status").eqInt("categoryId", "categoryId").build()
	if err != nil {
		respondBadRequest(c, err)
		return
	}
	items, total, err := h.svc.Query(c.Request.Context(), q)
	if err != nil {
		respondError(c, err)
		return
	}
	res := make([]postResponse, 0, len(items))
	for i := range items {
		res = append(res, newPostResponse(&items[i]))
	}
	setTotalCount(c, total)
	c.JSON(http.StatusOK, res)
}

func (h *PostHandler) get(c *gin.Context) {
	id, ok := parseID(c)
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

func (h *PostHandler) create(c *gin.Context) {
	var req postCreateRequest
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

func (h *PostHandler) update(c *gin.Context) {
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

	in := port.UpdatePostInput{
		Title:      cur.Title,
		Slug:       cur.Slug,
		Excerpt:    cur.Excerpt,
		Content:    cur.Content,
		Status:     cur.Status,
		CategoryID: cur.CategoryID,
		TagIDs:     nil, // nil = 태그 변경 없음
	}
	p.apply("title", &in.Title)
	p.apply("slug", &in.Slug)
	p.apply("excerpt", &in.Excerpt)
	p.apply("content", &in.Content)
	p.apply("status", &in.Status)
	p.apply("categoryId", &in.CategoryID) // null → 카테고리 해제
	if p.has("tagIds") {
		// 키가 있으면 교체(빈 배열·null 은 전체 해제). apply 가 null 을 nil 로 만들 수 있으므로
		// 다시 빈 슬라이스로 보정해 서비스가 '변경 없음(nil)'으로 오인하지 않게 한다.
		in.TagIDs = []uint{}
		p.apply("tagIds", &in.TagIDs)
		if in.TagIDs == nil {
			in.TagIDs = []uint{}
		}
	}
	if p.err != nil {
		respondBadRequest(c, p.err)
		return
	}

	out, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, newPostResponse(out))
}

func (h *PostHandler) delete(c *gin.Context) {
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
