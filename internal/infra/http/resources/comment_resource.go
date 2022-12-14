package resources

import (
	"nix_education/internal/domain"
)

type CommentDto struct {
	Id     int64  `json:"id"`
	PostId int64  `json:"post_id"`
	Name   string `json:"name"`
	Body   string `json:"body"`
}

type CommentsDto struct {
	Items []CommentDto `json:"items"`
	Total uint64       `json:"total"`
	Pages uint64       `json:"pages"`
}

func (d CommentDto) DomainToDto(comment domain.Comment) CommentDto {
	return CommentDto{
		Id:     comment.Id,
		PostId: comment.PostId,
		Name:   comment.Name,
		Body:   comment.Body,
	}
}

func (d CommentDto) DomainToDtoCollection(comments domain.Comments) CommentsDto {
	result := make([]CommentDto, len(comments.Items))

	for i := range comments.Items {
		result[i] = d.DomainToDto(comments.Items[i])
	}

	return CommentsDto{Items: result, Pages: comments.Pages, Total: comments.Total}
}
