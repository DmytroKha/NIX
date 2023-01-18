package requests

import (
	"nix_education/internal/infra/database"
)

type PostRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

func (r PostRequest) ToDomainModel() (database.Post, error) {
	var pst database.Post
	pst.Title = r.Title
	pst.Body = r.Body
	return pst, nil
}
