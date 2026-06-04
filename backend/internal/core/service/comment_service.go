package service

import (
	"context"
	"strings"

	"github.com/yoophi/refine-cms/backend/internal/core/domain"
	"github.com/yoophi/refine-cms/backend/internal/core/port"
)

type commentService struct {
	repo     port.CommentRepository
	postRepo port.PostRepository
}

// NewCommentService 는 CommentService 구현을 생성한다.
// 댓글 생성 시 대상 게시글 존재 여부를 확인하기 위해 PostRepository 도 주입받는다.
func NewCommentService(repo port.CommentRepository, postRepo port.PostRepository) port.CommentService {
	return &commentService{repo: repo, postRepo: postRepo}
}

func (s *commentService) Create(ctx context.Context, in port.CreateCommentInput) (*domain.Comment, error) {
	if in.PostID == 0 || strings.TrimSpace(in.AuthorName) == "" || strings.TrimSpace(in.Content) == "" {
		return nil, domain.ErrInvalidInput
	}
	// 대상 게시글이 존재하지 않으면 ErrNotFound 가 전파된다.
	if _, err := s.postRepo.GetByID(ctx, in.PostID); err != nil {
		return nil, err
	}

	status := in.Status
	if status == "" {
		status = domain.CommentStatusPending
	}
	if !status.Valid() {
		return nil, domain.ErrInvalidInput
	}
	c := &domain.Comment{
		PostID:      in.PostID,
		ParentID:    in.ParentID,
		UserID:      in.UserID,
		AuthorName:  in.AuthorName,
		AuthorEmail: in.AuthorEmail,
		Content:     in.Content,
		Status:      status,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *commentService) Get(ctx context.Context, id uint) (*domain.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *commentService) ListByPost(ctx context.Context, postID uint) ([]domain.Comment, error) {
	return s.repo.ListByPost(ctx, postID)
}

func (s *commentService) Query(ctx context.Context, q port.ListQuery) ([]domain.Comment, int, error) {
	return s.repo.Query(ctx, q)
}

func (s *commentService) Update(ctx context.Context, id uint, in port.UpdateCommentInput) (*domain.Comment, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.Content) == "" {
		return nil, domain.ErrInvalidInput
	}
	status := in.Status
	if status == "" {
		status = c.Status
	}
	if !status.Valid() {
		return nil, domain.ErrInvalidInput
	}
	c.Content = in.Content
	c.Status = status
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *commentService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// ownedComment 는 댓글을 조회하고 ownerID 소유 여부를 확인한다.
func (s *commentService) ownedComment(ctx context.Context, id, ownerID uint) (*domain.Comment, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c.UserID == nil || *c.UserID != ownerID {
		return nil, domain.ErrForbidden
	}
	return c, nil
}

func (s *commentService) UpdateOwnedContent(ctx context.Context, id, ownerID uint, content string) (*domain.Comment, error) {
	if strings.TrimSpace(content) == "" {
		return nil, domain.ErrInvalidInput
	}
	c, err := s.ownedComment(ctx, id, ownerID)
	if err != nil {
		return nil, err
	}
	c.Content = content // 상태(status)는 회원이 변경할 수 없다(관리자 전용)
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *commentService) DeleteOwned(ctx context.Context, id, ownerID uint) error {
	if _, err := s.ownedComment(ctx, id, ownerID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
