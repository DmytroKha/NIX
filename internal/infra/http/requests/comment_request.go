package requests

import "NIX/internal/domain"

type CommentRequest struct {
	PostId int64  `json:"postId"`
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func (r CommentRequest) ToDomainModel() (domain.Comment, error) {
	var cmt domain.Comment
	cmt.Id = r.Id
	cmt.PostId = r.PostId
	cmt.Name = r.Name
	cmt.Email = r.Email
	cmt.Body = r.Body

	return cmt, nil
}
