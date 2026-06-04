package domain

import "time"

// User 는 공개 웹앱의 회원이다(관리자 AdminUser 와 분리).
type User struct {
	ID           uint
	Email        string
	Name         string
	PasswordHash string // 응답에 절대 노출하지 않는다
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
