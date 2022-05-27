package dto

type Reply struct {
	CommentId string `validate:"required"`
	AuthorId  string `validate:"required"`
	Body      string `validate:"required"`
}
