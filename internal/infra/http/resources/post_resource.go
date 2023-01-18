package resources

import (
	"nix_education/internal/infra/database"
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

func (d PostDto) DatabaseToDto(post database.Post) PostDto {
	return PostDto{
		Id:     post.Id,
		UserId: post.UserId,
		Title:  post.Title,
		Body:   post.Body,
	}
}

func (d PostDto) DatabaseToDtoCollection(posts database.Posts) PostsDto {
	result := make([]PostDto, len(posts.Items))
	for i := range posts.Items {
		result[i] = d.DatabaseToDto(posts.Items[i])
	}
	return PostsDto{Items: result, Pages: posts.Pages, Total: posts.Total}
}
