package domain

import "time"

// Category 는 게시글을 분류하는 카테고리이다. 자기참조(ParentID)로 트리 구조를 표현할 수 있다.
type Category struct {
	ID          uint
	Name        string
	Slug        string
	Description string
	ParentID    *uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
