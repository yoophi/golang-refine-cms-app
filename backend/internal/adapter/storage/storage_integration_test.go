package storage

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/yoophi/refine-cms/backend/internal/config"
	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

// newTestDB 는 임시 파일 sqlite DB 를 열고 마이그레이션한 뒤 정리를 등록한다.
func newTestDB(t *testing.T) (*sqlx.DB, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	db, driver, err := NewDB(config.DBConfig{Driver: "sqlite", SQLitePath: path})
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	if err := Migrate(db, driver); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db, driver
}

func TestComment_UserIDRoundTrip(t *testing.T) {
	db, driver := newTestDB(t)
	ctx := context.Background()

	users := NewUserRepository(db, driver)
	posts := NewPostRepository(db, driver)
	comments := NewCommentRepository(db, driver)

	u := &domain.User{Email: "u@example.com", Name: "U", PasswordHash: "h"}
	if err := users.Create(ctx, u); err != nil {
		t.Fatalf("user create: %v", err)
	}
	p := &domain.Post{Title: "T", Slug: "t", Status: domain.PostStatusPublished}
	if err := posts.Create(ctx, p, nil); err != nil {
		t.Fatalf("post create: %v", err)
	}

	c := &domain.Comment{PostID: p.ID, UserID: &u.ID, AuthorName: "U", Content: "hi", Status: domain.CommentStatusPending}
	if err := comments.Create(ctx, c); err != nil {
		t.Fatalf("comment create: %v", err)
	}
	got, err := comments.GetByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("comment get: %v", err)
	}
	if got.UserID == nil || *got.UserID != u.ID {
		t.Fatalf("user_id round-trip 실패: got %v, want %d", got.UserID, u.ID)
	}
}

func TestCategoryQuery_LikeEscape(t *testing.T) {
	db, driver := newTestDB(t)
	ctx := context.Background()

	repo := NewCategoryRepository(db, driver)
	for _, name := range []string{"a_b", "axb"} {
		if err := repo.Create(ctx, &domain.Category{Name: name, Slug: name}); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}

	// LIKE 와일드카드 '_' 가 이스케이프되면 "a_b"(리터럴)만 매칭하고 "axb"는 제외돼야 한다.
	q := port.ListQuery{Filters: []port.Filter{{Field: "name", Op: port.OpLike, Value: "a_b"}}}
	items, total, err := repo.Query(ctx, q)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Name != "a_b" {
		t.Fatalf("LIKE escape 실패: total=%d items=%v (기대: a_b 1건)", total, items)
	}
}
