package storage

import (
	"strings"

	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// fieldMap 은 camelCase 입력 필드명을 실제 컬럼명으로 매핑하는 화이트리스트이다.
// 매핑에 없는 필드는 SQL 로 변환되지 않으므로 임의 컬럼 주입을 차단한다(안전).
type fieldMap map[string]string

// listClauses 는 ListQuery 로부터 생성된 SQL 조각이다.
type listClauses struct {
	where     string // " WHERE ..." 또는 ""
	whereArgs []any
	order     string // " ORDER BY ..."
	limit     string // " LIMIT ? OFFSET ?" 또는 ""
	limitArgs []any
}

// buildListClauses 는 화이트리스트(filterCols/sortCols)를 적용해 안전한 SQL 조각을 만든다.
// 플레이스홀더는 `?` 로 생성하며 호출 측에서 Rebind 한다.
func buildListClauses(q port.ListQuery, filterCols, sortCols fieldMap, defaultOrder string) listClauses {
	var c listClauses

	// --- WHERE ---
	var conds []string
	for _, f := range q.Filters {
		col, ok := filterCols[f.Field]
		if !ok {
			continue // 화이트리스트 외 필드 무시
		}
		switch f.Op {
		case port.OpEq:
			conds = append(conds, col+" = ?")
			c.whereArgs = append(c.whereArgs, f.Value)
		case port.OpNe:
			conds = append(conds, col+" <> ?")
			c.whereArgs = append(c.whereArgs, f.Value)
		case port.OpLike:
			// 대소문자 무시 부분 일치 (sqlite/postgres 공통).
			// 사용자 입력의 LIKE 와일드카드(%,_)와 이스케이프 문자(\)를 이스케이프해 와일드카드 주입 차단.
			conds = append(conds, "LOWER("+col+") LIKE LOWER(?) ESCAPE '\\'")
			c.whereArgs = append(c.whereArgs, "%"+escapeLike(toString(f.Value))+"%")
		case port.OpGt:
			conds = append(conds, col+" > ?")
			c.whereArgs = append(c.whereArgs, f.Value)
		case port.OpGte:
			conds = append(conds, col+" >= ?")
			c.whereArgs = append(c.whereArgs, f.Value)
		case port.OpLt:
			conds = append(conds, col+" < ?")
			c.whereArgs = append(c.whereArgs, f.Value)
		case port.OpLte:
			conds = append(conds, col+" <= ?")
			c.whereArgs = append(c.whereArgs, f.Value)
		}
	}
	if len(conds) > 0 {
		c.where = " WHERE " + strings.Join(conds, " AND ")
	}

	// --- ORDER BY ---
	var orders []string
	for _, s := range q.Sorts {
		col, ok := sortCols[s.Field]
		if !ok {
			continue
		}
		dir := "ASC"
		if s.Desc {
			dir = "DESC"
		}
		orders = append(orders, col+" "+dir)
	}
	if len(orders) > 0 {
		c.order = " ORDER BY " + strings.Join(orders, ", ")
	} else {
		c.order = " ORDER BY " + defaultOrder
	}

	// --- LIMIT / OFFSET ---
	if q.Limit > 0 {
		c.limit = " LIMIT ? OFFSET ?"
		c.limitArgs = []any{q.Limit, q.Offset}
	}

	return c
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// escapeLike 는 LIKE 패턴의 특수문자(\ % _)를 ESCAPE '\' 기준으로 이스케이프한다.
// 백슬래시를 먼저 처리해야 이중 이스케이프를 피한다.
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "%", `\%`)
	s = strings.ReplaceAll(s, "_", `\_`)
	return s
}
