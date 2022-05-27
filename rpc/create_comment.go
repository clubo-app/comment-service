package rpc

import (
	"context"

	"github.com/clubo-app/comment-service/dto"
	cg "github.com/clubo-app/protobuf/comment"
	"github.com/go-playground/validator/v10"
)

func (s commentServer) CreateComment(ctx context.Context, req *cg.CreateCommentRequest) (*cg.Comment, error) {
	dc := dto.Comment{
		PartyId:  req.PartyId,
		AuthorId: req.AuthorId,
		Body:     req.Body,
	}

	v := validator.New()
	err := v.Struct(dc)
	if err != nil {
		return nil, err
	}

	c, err := s.cs.Create(ctx, dc)
	if err != nil {
		return nil, err
	}

	return c.ToGRPCComment(), nil
}
