package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
	"github.com/yoophi/refine-cms/backend/internal/core/service"
)

func TestComment_Ownership(t *testing.T) {
	repo := newFakeCommentRepo()
	svc := service.NewCommentService(repo, newFakePostRepo(1))
	ctx := context.Background()

	owner := uint(10)
	other := uint(20)

	created, err := svc.Create(ctx, port.CreateCommentInput{
		PostID: 1, UserID: &owner, AuthorName: "A", Content: "원본",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.UserID == nil || *created.UserID != owner || created.Status != domain.CommentStatusPending {
		t.Fatalf("created mismatch: %+v", created)
	}

	// 타인 수정/삭제 → 403
	if _, err := svc.UpdateOwnedContent(ctx, created.ID, other, "해킹"); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("update by other: want ErrForbidden, got %v", err)
	}
	if err := svc.DeleteOwned(ctx, created.ID, other); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("delete by other: want ErrForbidden, got %v", err)
	}

	// 본인 수정 → 성공
	upd, err := svc.UpdateOwnedContent(ctx, created.ID, owner, "수정됨")
	if err != nil || upd.Content != "수정됨" {
		t.Fatalf("update by owner: %v, %+v", err, upd)
	}

	// 본인 삭제 → 성공
	if err := svc.DeleteOwned(ctx, created.ID, owner); err != nil {
		t.Fatalf("delete by owner: %v", err)
	}
	if _, err := svc.Get(ctx, created.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("after delete: want ErrNotFound, got %v", err)
	}
}

func TestComment_OwnershipEmptyContent(t *testing.T) {
	repo := newFakeCommentRepo()
	svc := service.NewCommentService(repo, newFakePostRepo(1))
	ctx := context.Background()
	owner := uint(10)
	c, _ := svc.Create(ctx, port.CreateCommentInput{PostID: 1, UserID: &owner, AuthorName: "A", Content: "x"})

	if _, err := svc.UpdateOwnedContent(ctx, c.ID, owner, "  "); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("empty content: want ErrInvalidInput, got %v", err)
	}
}

func TestComment_CreateMissingPost(t *testing.T) {
	svc := service.NewCommentService(newFakeCommentRepo(), newFakePostRepo()) // 게시글 없음
	owner := uint(10)
	if _, err := svc.Create(context.Background(), port.CreateCommentInput{PostID: 99, UserID: &owner, AuthorName: "A", Content: "x"}); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("create on missing post: want ErrNotFound, got %v", err)
	}
}
