package dto

type Comment struct {
	PartyId  string `validate:"required"`
	AuthorId string `validate:"required"`
	Body     string `validate:"required"`
}
