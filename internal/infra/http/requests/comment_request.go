package requests

import (
	"nix_education/internal/infra/database"
)

type CommentRequest struct {
	Name string `json:"name" validate:"required"`
	Body string `json:"body" validate:"required"`
}

func (r CommentRequest) ToDatabaseModel() (database.Comment, error) {
	var cmt database.Comment
	cmt.Name = r.Name
	cmt.Body = r.Body

	return cmt, nil
}
