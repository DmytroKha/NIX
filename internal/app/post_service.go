package app

import (
	"NIX/internal/domain"
	"NIX/internal/infra/database"
	"log"
)

type PostService interface {
	Save(post domain.Post) (domain.Post, error)
	Update(post domain.Post) (domain.Post, error)
	Find(id int64) (domain.Post, error)
	FindAll(userId int64, p domain.Pagination) (domain.Posts, error)
	Delete(id int64) error
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

func (s postService) FindAll(userId int64, p domain.Pagination) (domain.Posts, error) {
	posts, err := s.postRepo.FindAll(userId, p)
	if err != nil {
		log.Print(err)
		return domain.Posts{}, err
	}

	return posts, nil
}

func (s postService) Update(p domain.Post) (domain.Post, error) {
	post, err := s.postRepo.Update(p)
	if err != nil {
		log.Print(err)
		return domain.Post{}, err
	}

	return post, nil
}

func (s postService) Delete(id int64) error {
	err := s.postRepo.Delete(id)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
