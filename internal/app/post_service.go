package app

import (
	"errors"
	"log"
	"nix_education/internal/domain"
	"nix_education/internal/infra/database"
	"nix_education/internal/infra/http/resources"
)

//go:generate mockery --dir . --name PostService --output ./mocks
type PostService interface {
	Save(post database.Post) (resources.PostDto, error)
	Update(post database.Post) (resources.PostDto, error)
	Find(id int64) (resources.PostDto, error)
	FindAll(p domain.Pagination) (resources.PostsDto, error)
	Delete(id, userId int64) error
}

type postService struct {
	postRepo database.PostRepository
}

func NewPostService(r database.PostRepository) PostService {
	return postService{
		postRepo: r,
	}
}

func (s postService) Save(p database.Post) (resources.PostDto, error) {
	post, err := s.postRepo.Save(p)
	if err != nil {
		log.Print(err)
		return resources.PostDto{}, err
	}
	var postDto resources.PostDto
	return postDto.DatabaseToDto(post), nil
}

func (s postService) Find(id int64) (resources.PostDto, error) {
	post, err := s.postRepo.Find(id)
	if err != nil {
		log.Print(err)
		return resources.PostDto{}, err
	}
	var postDto resources.PostDto
	return postDto.DatabaseToDto(post), nil
}

func (s postService) FindAll(p domain.Pagination) (resources.PostsDto, error) {
	posts, err := s.postRepo.FindAll(p)
	if err != nil {
		log.Print(err)
		return resources.PostsDto{}, err
	}
	var postDto resources.PostDto
	return postDto.DatabaseToDtoCollection(posts), nil
}

func (s postService) Update(p database.Post) (resources.PostDto, error) {
	findPost, err := s.Find(p.Id)
	if err != nil {
		log.Print(err)
		return resources.PostDto{}, err
	}
	if findPost.UserId != p.UserId {
		err := errors.New("user id mismatch")
		log.Print(err)
		return resources.PostDto{}, err
	}
	post, err := s.postRepo.Update(p)
	if err != nil {
		log.Print(err)
		return resources.PostDto{}, err
	}
	var postDto resources.PostDto
	return postDto.DatabaseToDto(post), nil
}

func (s postService) Delete(id, userId int64) error {
	findPost, err := s.Find(id)
	if err != nil {
		log.Print(err)
		return err
	}
	if findPost.UserId != userId {
		err = errors.New("user id mismatch")
		log.Print(err)
		return err
	}
	err = s.postRepo.Delete(id)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
