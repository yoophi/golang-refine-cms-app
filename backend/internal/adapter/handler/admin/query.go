package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

const (
	totalCountHeader = "X-Total-Count"
	// maxPageSize 는 한 번에 조회할 수 있는 최대 행 수(무제한 조회/DoS 방지).
	maxPageSize = 200
)

// listQuery 는 refine simple-rest 쿼리 파라미터를 안전하게 ListQuery 로 파싱하는 빌더이다.
// 정수 필터 파싱 실패 등 입력 오류를 누적해 핸들러가 400 으로 응답할 수 있게 한다.
type listQuery struct {
	c   *gin.Context
	q   port.ListQuery
	err error
}

func newListQuery(c *gin.Context) *listQuery {
	offset, limit := parsePaging(c)
	return &listQuery{c: c, q: port.ListQuery{Offset: offset, Limit: limit, Sorts: parseSort(c)}}
}

// like 는 부분 일치(LIKE) 문자열 필터를 추가한다(값이 있으면).
func (b *listQuery) like(param, field string) *listQuery {
	if v := b.c.Query(param); v != "" {
		b.q.Filters = append(b.q.Filters, port.Filter{Field: field, Op: port.OpLike, Value: v})
	}
	return b
}

// eqStr 은 정확 일치 문자열 필터를 추가한다.
func (b *listQuery) eqStr(param, field string) *listQuery {
	if v := b.c.Query(param); v != "" {
		b.q.Filters = append(b.q.Filters, port.Filter{Field: field, Op: port.OpEq, Value: v})
	}
	return b
}

// eqInt 는 정수 정확 일치 필터를 추가한다. 값이 정수가 아니면 400 에러를 누적한다
// (조용히 무시해 전체 조회로 빠지지 않도록).
func (b *listQuery) eqInt(param, field string) *listQuery {
	v := b.c.Query(param)
	if v == "" {
		return b
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		if b.err == nil {
			b.err = fmt.Errorf("%s 는 정수여야 합니다", param)
		}
		return b
	}
	b.q.Filters = append(b.q.Filters, port.Filter{Field: field, Op: port.OpEq, Value: n})
	return b
}

// build 는 완성된 ListQuery 와 누적된 입력 오류를 반환한다.
func (b *listQuery) build() (port.ListQuery, error) { return b.q, b.err }

// parsePaging 은 refine 의 _start/_end 를 offset/limit 으로 변환한다.
// 누락/오입력이거나 범위가 비정상이면 무제한 조회 대신 maxPageSize 로 제한한다.
func parsePaging(c *gin.Context) (offset, limit int) {
	start, _ := strconv.Atoi(c.Query("_start"))
	end, _ := strconv.Atoi(c.Query("_end"))
	if start < 0 {
		start = 0
	}
	limit = end - start
	if limit <= 0 || limit > maxPageSize {
		limit = maxPageSize
	}
	return start, limit
}

// parseSort 은 _sort(콤마 구분 필드)/_order(콤마 구분 asc|desc)를 SortField 로 변환한다.
// 필드명은 camelCase 그대로 두며 리포지토리가 화이트리스트로 검증한다.
func parseSort(c *gin.Context) []port.SortField {
	sortRaw := strings.TrimSpace(c.Query("_sort"))
	if sortRaw == "" {
		return nil
	}
	fields := strings.Split(sortRaw, ",")
	orders := strings.Split(c.Query("_order"), ",")

	sorts := make([]port.SortField, 0, len(fields))
	for i, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		desc := i < len(orders) && strings.EqualFold(strings.TrimSpace(orders[i]), "desc")
		sorts = append(sorts, port.SortField{Field: f, Desc: desc})
	}
	return sorts
}

// setTotalCount 는 refine 페이지네이션에 필수인 X-Total-Count 헤더를 설정한다.
func setTotalCount(c *gin.Context, total int) {
	c.Header(totalCountHeader, strconv.Itoa(total))
}
