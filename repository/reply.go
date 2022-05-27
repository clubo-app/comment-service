package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/clubo-app/comment-service/datastruct"
	"github.com/go-playground/validator/v10"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

const (
	REPLIES = "replies"
)

var repliesMetadata = table.Metadata{
	Name:    REPLIES,
	Columns: []string{"id", "comment_id", "author_id", "body"},
	PartKey: []string{"comment_id", "id"},
}

type ReplyRepository interface {
	Create(ctx context.Context, p datastruct.Reply) (datastruct.Reply, error)
	Delete(ctx context.Context, uId, cId, rId string) error
	GetByComment(ctx context.Context, cId string, page []byte, limit uint32) ([]datastruct.Reply, []byte, error)
}

type replyRepository struct {
	sess *gocqlx.Session
}

func (r *replyRepository) Create(ctx context.Context, rp datastruct.Reply) (datastruct.Reply, error) {
	v := validator.New()
	err := v.Struct(rp)
	if err != nil {
		return datastruct.Reply{}, err
	}

	stmt, names := qb.
		Insert(REPLIES).
		Columns(repliesMetadata.Columns...).
		ToCql()

	err = r.sess.
		Query(stmt, names).
		BindStruct(rp).
		ExecRelease()
	if err != nil {
		return datastruct.Reply{}, err
	}

	return rp, err
}

func (r *replyRepository) Delete(ctx context.Context, uId, cId, rId string) error {
	stmt, names := qb.
		Delete(REPLIES).
		Where(qb.Eq("id")).
		Where(qb.Eq("id")).
		If(qb.Eq("author_id")).
		ToCql()

	err := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"id": rId, "comment_id": cId, "author_id": uId})).
		ExecRelease()
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return errors.New("you can only Delete your own Replies")
		}
		return err
	}
	return nil
}

func (r *replyRepository) GetByComment(ctx context.Context, cId string, page []byte, limit uint32) ([]datastruct.Reply, []byte, error) {
	var result []datastruct.Reply
	stmt, names := qb.
		Select(REPLIES).
		Where(qb.Eq("comment_id")).
		ToCql()

	q := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"comment_id": cId}))
	defer q.Release()

	q.PageState(page)
	if limit == 0 {
		q.PageSize(10)
	} else {
		q.PageSize(int(limit))
	}

	iter := q.Iter()
	err := iter.Select(&result)
	if err != nil {
		return []datastruct.Reply{}, nil, errors.New("no replies found")
	}

	return result, iter.PageState(), nil

}
