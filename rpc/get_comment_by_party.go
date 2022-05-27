package rpc

import (
	"context"
	"encoding/base64"

	"github.com/clubo-app/packages/utils"
	cg "github.com/clubo-app/protobuf/comment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s commentServer) GetCommentByParty(ctx context.Context, req *cg.GetByPartyRequest) (*cg.PagedComments, error) {
	p, err := base64.URLEncoding.DecodeString(req.NextPage)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Next Page Param")
	}

	cs, p, err := s.cs.GetByParty(ctx, req.PartyId, p, req.Limit)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	nextPage := base64.URLEncoding.EncodeToString(p)

	pc := make([]*cg.Comment, len(cs))
	for i, c := range cs {
		pc[i] = c.ToGRPCComment()
	}

	return &cg.PagedComments{Comments: pc, NextPage: nextPage}, nil
}
