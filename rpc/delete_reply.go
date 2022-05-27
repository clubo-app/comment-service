package rpc

import (
	"context"

	cg "github.com/clubo-app/protobuf/comment"
	common "github.com/clubo-app/protobuf/common"
)

func (s commentServer) DeleteReply(ctx context.Context, req *cg.DeleteReplyRequest) (*common.MessageResponse, error) {
	err := s.rs.Delete(ctx, req.AuthorId, req.CommentId, req.ReplyId)
	if err != nil {
		return nil, err
	}

	return &common.MessageResponse{Message: "Reply removed"}, nil
}
