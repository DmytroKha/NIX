package requests

import "NIX/internal/domain"

type PostRequest struct {
	//UserId int64 `json:"userId,omitempty"`
	//Id     int64  `json:"id,omitempty"`
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

func (r PostRequest) ToDomainModel() (domain.Post, error) {

	var pst domain.Post
	//pst.Id = r.Id
	//pst.UserId = r.UserId
	pst.Title = r.Title
	pst.Body = r.Body

	return pst, nil
}
