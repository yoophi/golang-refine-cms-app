package domain

import "time"

// CommentStatus 는 댓글의 검수 상태를 나타낸다.
type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"
	CommentStatusApproved CommentStatus = "approved"
	CommentStatusSpam     CommentStatus = "spam"
)

// Valid 는 정의된 상태값인지 검증한다.
func (s CommentStatus) Valid() bool {
	switch s {
	case CommentStatusPending, CommentStatusApproved, CommentStatusSpam:
		return true
	default:
		return false
	}
}

// Comment 는 게시글에 달리는 댓글이다. ParentID 로 대댓글(스레드)을 표현한다.
type Comment struct {
	ID          uint
	PostID      uint
	ParentID    *uint
	AuthorName  string
	AuthorEmail string
	Content     string
	Status      CommentStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
