package rpc

import (
	"context"
	"encoding/base64"

	"github.com/clubo-app/packages/utils"
	cg "github.com/clubo-app/protobuf/comment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s commentServer) GetReplyByComment(ctx context.Context, req *cg.GetReplyByCommentRequest) (*cg.PagedReply, error) {
	p, err := base64.URLEncoding.DecodeString(req.NextPage)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Next Page Param")
	}

	rs, p, err := s.rs.GetByComment(ctx, req.CommentId, p, req.Limit)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	nextPage := base64.URLEncoding.EncodeToString(p)

	pr := make([]*cg.Reply, len(rs))
	for i, r := range rs {
		pr[i] = r.ToGRPCReply()
	}

	return &cg.PagedReply{Replies: pr, NextPage: nextPage}, nil
}
