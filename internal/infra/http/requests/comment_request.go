package requests

import "NIX/internal/domain"

type CommentRequest struct {
	//PostId int64  `json:"postId,omitempty"`
	//Id     int64  `json:"id,omitempty"`
	Name string `json:"name" validate:"required"`
	//Email string `json:"email"`
	Body string `json:"body" validate:"required"`
}

func (r CommentRequest) ToDomainModel() (domain.Comment, error) {
	var cmt domain.Comment
	//cmt.Id = r.Id
	//cmt.PostId = r.PostId
	cmt.Name = r.Name
	//cmt.Email = r.Email
	cmt.Body = r.Body

	return cmt, nil
}
