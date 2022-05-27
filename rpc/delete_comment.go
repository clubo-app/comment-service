package rpc

import (
	"context"

	cg "github.com/clubo-app/protobuf/comment"
	common "github.com/clubo-app/protobuf/common"
)

func (s commentServer) DeleteComment(ctx context.Context, req *cg.DeleteCommentRequest) (*common.MessageResponse, error) {
	err := s.cs.Delete(ctx, req.AuthorId, req.PartyId, req.CommentId)
	if err != nil {
		return nil, err
	}

	return &common.MessageResponse{Message: "Comment removed"}, nil
}
