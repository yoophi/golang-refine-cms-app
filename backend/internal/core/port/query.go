package port

// 목록 조회(refine simple-rest)용 공통 쿼리 모델.
// 인바운드 어댑터가 HTTP 쿼리 파라미터를 ListQuery 로 파싱하고,
// 아웃바운드 어댑터(리포지토리)가 화이트리스트로 안전하게 SQL 로 변환한다.

// FilterOp 는 목록 필터 연산자이다.
type FilterOp string

const (
	OpEq   FilterOp = "eq"   // {field}={value}
	OpNe   FilterOp = "ne"   // {field}_ne=
	OpLike FilterOp = "like" // {field}_like= (부분 일치)
	OpGt   FilterOp = "gt"   // {field}_gt=
	OpGte  FilterOp = "gte"  // {field}_gte=
	OpLt   FilterOp = "lt"   // {field}_lt=
	OpLte  FilterOp = "lte"  // {field}_lte=
)

// Filter 는 단일 필드 필터 조건이다.
// Field 는 camelCase(어댑터 입력)이며, 리포지토리가 화이트리스트로 컬럼에 매핑한다.
// Value 는 컬럼 타입에 맞게(정수/문자열) 어댑터에서 변환되어 들어온다.
type Filter struct {
	Field string
	Op    FilterOp
	Value any
}

// SortField 는 정렬 조건이다(Field 는 camelCase).
type SortField struct {
	Field string
	Desc  bool
}

// ListQuery 는 목록 조회 파라미터이다.
type ListQuery struct {
	Offset  int // _start
	Limit   int // _end - _start (0 이면 전체)
	Sorts   []SortField
	Filters []Filter
}
