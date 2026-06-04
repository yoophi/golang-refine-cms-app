package admin

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

const totalCountHeader = "X-Total-Count"

// parsePaging 은 refine 의 _start/_end 를 offset/limit 으로 변환한다.
// 둘 다 유효하지 않으면 limit=0(전체)을 반환한다.
func parsePaging(c *gin.Context) (offset, limit int) {
	start, errS := strconv.Atoi(c.Query("_start"))
	end, errE := strconv.Atoi(c.Query("_end"))
	if errS != nil || errE != nil || start < 0 || end <= start {
		return 0, 0
	}
	return start, end - start
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

// baseQuery 는 페이지네이션 + 정렬을 채운 ListQuery 를 만든다(필터는 핸들러가 리소스별로 추가).
func baseQuery(c *gin.Context) port.ListQuery {
	offset, limit := parsePaging(c)
	return port.ListQuery{
		Offset: offset,
		Limit:  limit,
		Sorts:  parseSort(c),
	}
}

// strFilter 는 쿼리 파라미터가 있으면 해당 연산자의 문자열 필터를 추가한다.
func strFilter(c *gin.Context, q *port.ListQuery, param, field string, op port.FilterOp) {
	if v := c.Query(param); v != "" {
		q.Filters = append(q.Filters, port.Filter{Field: field, Op: op, Value: v})
	}
}

// intFilter 는 쿼리 파라미터가 정수면 eq 필터를 추가한다(타입 안전: postgres 대비 int 로 변환).
func intFilter(c *gin.Context, q *port.ListQuery, param, field string) {
	if v := c.Query(param); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			q.Filters = append(q.Filters, port.Filter{Field: field, Op: port.OpEq, Value: n})
		}
	}
}

// setTotalCount 는 refine 페이지네이션에 필수인 X-Total-Count 헤더를 설정한다.
func setTotalCount(c *gin.Context, total int) {
	c.Header(totalCountHeader, strconv.Itoa(total))
}
