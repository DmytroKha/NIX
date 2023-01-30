package app

import (
	"errors"
	"log"
	"nix_education/internal/domain"
	"nix_education/internal/infra/database"
	"nix_education/internal/infra/http/requests"
)

//go:generate mockery --dir . --name PostService --output ./mocks
type PostService interface {
	Save(pst requests.PostRequest, userId int64) (database.Post, error)
	Update(pst requests.PostRequest, id, userId int64) (database.Post, error)
	Find(id int64) (database.Post, error)
	FindAll(p domain.Pagination) (database.Posts, error)
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

func (s postService) Save(pst requests.PostRequest, userId int64) (database.Post, error) {
	p, err := pst.ToDatabaseModel()
	if err != nil {
		log.Print(err)
		return database.Post{}, err
	}

	p.UserId = userId
	post, err := s.postRepo.Save(p)
	if err != nil {
		log.Print(err)
		return database.Post{}, err
	}
	return post, nil
}

func (s postService) Find(id int64) (database.Post, error) {
	post, err := s.postRepo.Find(id)
	if err != nil {
		log.Print(err)
		return database.Post{}, err
	}
	return post, nil
}

func (s postService) FindAll(p domain.Pagination) (database.Posts, error) {
	posts, err := s.postRepo.FindAll(p)
	if err != nil {
		log.Print(err)
		return database.Posts{}, err
	}
	return posts, nil
}

func (s postService) Update(pst requests.PostRequest, id, userId int64) (database.Post, error) {
	p, err := pst.ToDatabaseModel()
	if err != nil {
		log.Print(err)
		return database.Post{}, err
	}
	p.UserId = userId
	p.Id = id

	findPost, err := s.Find(id)
	if err != nil {
		log.Print(err)
		return database.Post{}, err
	}
	if findPost.UserId != p.UserId {
		err := errors.New("user id mismatch")
		log.Print(err)
		return database.Post{}, err
	}
	post, err := s.postRepo.Update(p)
	if err != nil {
		log.Print(err)
		return database.Post{}, err
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
