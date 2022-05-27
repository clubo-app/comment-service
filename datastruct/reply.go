package datastruct

import (
	"time"

	cg "github.com/clubo-app/protobuf/comment"
	"github.com/gofrs/uuid"
)

type Reply struct {
	Id        string `json:"id"         db:"id"         validate:"required"`
	CommentId string `json:"comment_id" db:"comment_id" validate:"required"`
	AuthorId  string `json:"author_id"  db:"author_id"  validate:"required"`
	Body      string `json:"body"       db:"body"       validate:"required"`
}

func (r Reply) ToGRPCReply() *cg.Reply {
	uuidv1, err := uuid.FromString(r.Id)
	if err != nil {
		return &cg.Reply{}
	}
	timestamp, err := uuid.TimestampFromV1(uuidv1)
	if err != nil {
		return &cg.Reply{}
	}
	t, err := timestamp.Time()
	if err != nil {
		return &cg.Reply{}
	}

	return &cg.Reply{
		Id:        r.Id,
		CommentId: r.CommentId,
		AuthorId:  r.AuthorId,
		Body:      r.Body,
		CreatedAt: t.UTC().Format(time.RFC3339),
	}
}
