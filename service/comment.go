package service

import (
	"context"
	"errors"

	"github.com/clubo-app/comment-service/datastruct"
	"github.com/clubo-app/comment-service/dto"
	"github.com/clubo-app/comment-service/repository"
	"github.com/gofrs/uuid"
)

type CommentService interface {
	Create(ctx context.Context, c dto.Comment) (datastruct.Comment, error)
	Delete(ctx context.Context, uId, pId, cId string) error
	GetByParty(ctx context.Context, pId string, page []byte, limit uint32) ([]datastruct.Comment, []byte, error)
	GetByPartyUser(ctx context.Context, pId, uId string) ([]datastruct.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (cs commentService) Create(ctx context.Context, c dto.Comment) (datastruct.Comment, error) {
	uuid, err := uuid.NewV1()
	if err != nil {
		return datastruct.Comment{}, errors.New("failed generate Comment id")
	}

	dc := datastruct.Comment{
		Id:       uuid.String(),
		PartyId:  c.PartyId,
		AuthorId: c.AuthorId,
		Body:     c.Body,
	}
	return cs.repo.Create(ctx, dc)
}

func (cs commentService) Delete(ctx context.Context, uId, pId, cId string) error {
	return cs.repo.Delete(ctx, uId, pId, cId)
}

func (cs commentService) GetByParty(ctx context.Context, pId string, page []byte, limit uint32) ([]datastruct.Comment, []byte, error) {
	return cs.repo.GetByParty(ctx, pId, page, limit)
}

func (cs commentService) GetByPartyUser(ctx context.Context, pId, uId string) ([]datastruct.Comment, error) {
	return cs.repo.GetByPartyUser(ctx, pId, uId)
}
