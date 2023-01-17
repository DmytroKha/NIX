package app

import (
	"errors"
	"log"
	"nix_education/internal/domain"
	"nix_education/internal/infra/database"
)

//go:generate mockery --dir . --name PostService --output ./mocks
type PostService interface {
	Save(post domain.Post) (domain.Post, error)
	Update(post domain.Post) (domain.Post, error)
	Find(id int64) (domain.Post, error)
	FindAll(p domain.Pagination) (domain.Posts, error)
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

func (s postService) Save(p domain.Post) (domain.Post, error) {
	post, err := s.postRepo.Save(p)
	if err != nil {
		log.Print(err)
		return domain.Post{}, err
	}

	return post, nil
}

func (s postService) Find(id int64) (domain.Post, error) {
	post, err := s.postRepo.Find(id)
	if err != nil {
		log.Print(err)
		return domain.Post{}, err
	}

	return post, nil
}

func (s postService) FindAll(p domain.Pagination) (domain.Posts, error) {
	posts, err := s.postRepo.FindAll(p)
	if err != nil {
		log.Print(err)
		return domain.Posts{}, err
	}

	return posts, nil
}

func (s postService) Update(p domain.Post) (domain.Post, error) {

	findPost, err := s.Find(p.Id)

	if err != nil {
		log.Print(err)
		return domain.Post{}, err
	}

	if findPost.UserId != p.UserId {
		err := errors.New("user id mismatch")
		log.Print(err)
		return domain.Post{}, err
	}

	post, err := s.postRepo.Update(p)
	if err != nil {
		log.Print(err)
		return domain.Post{}, err
	}

	return post, nil
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
