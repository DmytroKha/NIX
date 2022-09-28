package requests

import "NIX/internal/domain"

type PostRequest struct {
	UserId int64  `json:"userId"`
	Id     int64  `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (r PostRequest) ToDomainModel() (domain.Post, error) {

	var pst domain.Post
	pst.Id = r.Id
	pst.UserId = r.UserId
	pst.Title = r.Title
	pst.Body = r.Body

	return pst, nil
}
