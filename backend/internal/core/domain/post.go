package domain

import "time"

// PostStatus 는 게시글의 발행 상태를 나타낸다.
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// Valid 는 정의된 상태값인지 검증한다.
func (s PostStatus) Valid() bool {
	switch s {
	case PostStatusDraft, PostStatusPublished, PostStatusArchived:
		return true
	default:
		return false
	}
}

// Post 는 CMS 의 핵심 콘텐츠 엔티티이다.
type Post struct {
	ID          uint
	Title       string
	Slug        string
	Excerpt     string
	Content     string
	Status      PostStatus
	CategoryID  *uint
	Category    *Category
	Tags        []Tag
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
