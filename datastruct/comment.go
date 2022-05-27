package datastruct

import (
	"time"

	cg "github.com/clubo-app/protobuf/comment"
	"github.com/gofrs/uuid"
)

type Comment struct {
	Id       string `json:"id"        db:"id"         validate:"required"`
	PartyId  string `json:"party_id"  db:"party_id"   validate:"required"`
	AuthorId string `json:"author_id" db:"author_id"  validate:"required"`
	Body     string `json:"body"      db:"body"       validate:"required"`
}

func (c Comment) ToGRPCComment() *cg.Comment {
	uuidv1, err := uuid.FromString(c.Id)
	if err != nil {
		return &cg.Comment{}
	}
	timestamp, err := uuid.TimestampFromV1(uuidv1)
	if err != nil {
		return &cg.Comment{}
	}
	t, err := timestamp.Time()
	if err != nil {
		return &cg.Comment{}
	}

	return &cg.Comment{
		Id:        c.Id,
		PartyId:   c.PartyId,
		AuthorId:  c.AuthorId,
		Body:      c.Body,
		CreatedAt: t.UTC().Format(time.RFC3339),
	}
}
