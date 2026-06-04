package domain

import "time"

// Tag 는 게시글에 부착하는 라벨이다. 게시글과 다대다 관계를 가진다.
type Tag struct {
	ID        uint
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
