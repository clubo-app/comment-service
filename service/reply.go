package service

import (
	"context"
	"errors"

	"github.com/clubo-app/comment-service/datastruct"
	"github.com/clubo-app/comment-service/dto"
	"github.com/clubo-app/comment-service/repository"
	"github.com/gofrs/uuid"
)

type ReplyService interface {
	Create(ctx context.Context, p dto.Reply) (datastruct.Reply, error)
	Delete(ctx context.Context, uId, cId, rId string) error
	GetByComment(ctx context.Context, cId string, page []byte, limit uint32) ([]datastruct.Reply, []byte, error)
}

type replyService struct {
	repo repository.ReplyRepository
}

func NewReplyService(repo repository.ReplyRepository) ReplyService {
	return &replyService{repo: repo}
}

func (s replyService) Create(ctx context.Context, r dto.Reply) (datastruct.Reply, error) {
	uuid, err := uuid.NewV1()
	if err != nil {
		return datastruct.Reply{}, errors.New("failed generate Reply id")
	}

	dc := datastruct.Reply{
		Id:        uuid.String(),
		CommentId: r.CommentId,
		AuthorId:  r.AuthorId,
		Body:      r.Body,
	}
	return s.repo.Create(ctx, dc)
}

func (s replyService) Delete(ctx context.Context, uId, cId, rId string) error {
	return s.repo.Delete(ctx, uId, cId, rId)
}
func (s replyService) GetByComment(ctx context.Context, cId string, page []byte, limit uint32) ([]datastruct.Reply, []byte, error) {
	return s.repo.GetByComment(ctx, cId, page, limit)
}
