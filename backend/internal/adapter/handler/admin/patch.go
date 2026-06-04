package admin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// patcher 는 PATCH(부분 수정) 요청 본문을 키 존재 여부 기반으로 적용한다.
// refine 은 보통 전체 폼을 보내지만, 일부 필드만 보내는 부분 수정도 지원한다.
// 키가 존재하면(값이 null 이어도) 적용하고, 없으면 기존 값을 유지한다.
type patcher struct {
	raw map[string]json.RawMessage
	err error
}

// newPatcher 는 요청 본문을 키-원시값 맵으로 읽는다.
func newPatcher(c *gin.Context) (*patcher, error) {
	var raw map[string]json.RawMessage
	if err := c.ShouldBindJSON(&raw); err != nil {
		return nil, err
	}
	return &patcher{raw: raw}, nil
}

// apply 는 key 가 존재하면 dst(포인터)에 값을 디코드한다. 첫 에러를 보관한다.
func (p *patcher) apply(key string, dst any) {
	if p.err != nil {
		return
	}
	if v, ok := p.raw[key]; ok {
		p.err = json.Unmarshal(v, dst)
	}
}

// has 는 본문에 해당 key 가 존재하는지 반환한다.
func (p *patcher) has(key string) bool {
	_, ok := p.raw[key]
	return ok
}
