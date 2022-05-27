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
	COMMENTS_BY_PARTY = "comments_by_party"
	COMMENTS_BY_USER  = "comments_by_user"
)

var commentsMetadata = table.Metadata{
	Name:    COMMENTS_BY_PARTY,
	Columns: []string{"id", "party_id", "author_id", "body"},
	PartKey: []string{"party_id", "id"},
}

type CommentRepository interface {
	Create(ctx context.Context, p datastruct.Comment) (datastruct.Comment, error)
	Delete(ctx context.Context, uId, pId, cId string) error
	GetByParty(ctx context.Context, pId string, page []byte, limit uint32) ([]datastruct.Comment, []byte, error)
	GetByPartyUser(ctx context.Context, pId, uId string) ([]datastruct.Comment, error)
}

type commentRepository struct {
	sess *gocqlx.Session
}

func (r *commentRepository) Create(ctx context.Context, c datastruct.Comment) (datastruct.Comment, error) {
	v := validator.New()
	err := v.Struct(c)
	if err != nil {
		return datastruct.Comment{}, err
	}

	stmt, names := qb.
		Insert(COMMENTS_BY_PARTY).
		Columns(commentsMetadata.Columns...).
		ToCql()

	err = r.sess.
		Query(stmt, names).
		BindStruct(c).
		ExecRelease()
	if err != nil {
		return datastruct.Comment{}, err
	}

	return c, err
}

// https://github.com/scylladb/scylla/issues/10171
// TODO: currently deletion by index is not supported if supported create GSI on comment_id and delete by it
func (r *commentRepository) Delete(ctx context.Context, uId, pId, cId string) error {
	stmt, names := qb.
		Delete(COMMENTS_BY_PARTY).
		Where(qb.Eq("id")).
		Where(qb.Eq("party_id")).
		If(qb.Eq("author_id")).
		ToCql()

	err := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"id": cId, "party_id": pId, "author_id": uId})).
		ExecRelease()
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return errors.New("you can only Delete your own Comments")
		}
		return err
	}
	return nil
}

func (r *commentRepository) GetByParty(ctx context.Context, pId string, page []byte, limit uint32) ([]datastruct.Comment, []byte, error) {
	var result []datastruct.Comment
	stmt, names := qb.
		Select(COMMENTS_BY_PARTY).
		Where(qb.Eq("party_id")).
		ToCql()

	q := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"party_id": pId}))
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
		return []datastruct.Comment{}, nil, errors.New("no comments found")
	}

	return result, iter.PageState(), nil
}

func (r *commentRepository) GetByPartyUser(ctx context.Context, pId, uId string) ([]datastruct.Comment, error) {
	var result []datastruct.Comment
	stmt, names := qb.
		Select(COMMENTS_BY_USER).
		Where(qb.Eq("party_id")).
		Where(qb.Eq("author_id")).
		OrderBy("created_at", qb.ASC).
		ToCql()

	err := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"party_id": pId, "author_id": uId})).
		PageSize(10).
		Iter().
		Select(&result)
	if err != nil {
		return []datastruct.Comment{}, errors.New("no comments found")
	}

	return result, nil
}
