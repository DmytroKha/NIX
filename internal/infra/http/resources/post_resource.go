package resources

import (
	"NIX/internal/domain"
)

type PostDto struct {
	Id     int64  `json:"id"`
	UserId int64  `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type PostsDto struct {
	Items []PostDto `json:"items"`
	Total uint64    `json:"total"`
	Pages uint64    `json:"pages"`
}

func (d PostDto) DomainToDto(post domain.Post) PostDto {
	return PostDto{
		Id:     post.Id,
		UserId: post.UserId,
		Title:  post.Title,
		Body:   post.Body,
	}
}

func (d PostDto) DomainToDtoCollection(posts domain.Posts) PostsDto {
	result := make([]PostDto, len(posts.Items))

	for i := range posts.Items {
		result[i] = d.DomainToDto(posts.Items[i])
	}

	return PostsDto{Items: result, Pages: posts.Pages, Total: posts.Total}
}
