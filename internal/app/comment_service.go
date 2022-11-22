package app

import (
	"NIX/internal/domain"
	"NIX/internal/infra/database"
	"errors"
	"log"
)

type CommentService interface {
	Save(comment domain.Comment) (domain.Comment, error)
	Update(comment domain.Comment) (domain.Comment, error)
	Find(postId, id int64) (domain.Comment, error)
	FindAll(postId int64, p domain.Pagination) (domain.Comments, error)
	Delete(postId, id int64, email string) error
}

type commentService struct {
	commentRepo database.CommentRepository
	postServise PostService
}

func NewCommentService(r database.CommentRepository, ps PostService) CommentService {
	return commentService{
		commentRepo: r,
		postServise: ps,
	}
}

func (s commentService) Save(p domain.Comment) (domain.Comment, error) {

	_, err := s.postServise.Find(p.PostId)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}

	comment, err := s.commentRepo.Save(p)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}

	return comment, nil
}

func (s commentService) Find(postId, id int64) (domain.Comment, error) {
	comment, err := s.commentRepo.Find(postId, id)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}

	return comment, nil
}

func (s commentService) FindAll(postId int64, p domain.Pagination) (domain.Comments, error) {
	comments, err := s.commentRepo.FindAll(postId, p)
	if err != nil {
		log.Print(err)
		return domain.Comments{}, err
	}

	return comments, nil
}

func (s commentService) Update(p domain.Comment) (domain.Comment, error) {
	post, err := s.commentRepo.Find(p.PostId, p.Id)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}

	if p.Email != post.Email {
		err = errors.New("user email mismatch")
		return domain.Comment{}, err
	}

	comment, err := s.commentRepo.Update(p)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}

	return comment, nil
}

func (s commentService) Delete(postId, id int64, email string) error {
	post, err := s.commentRepo.Find(postId, id)
	if err != nil {
		log.Print(err)
		return err
	}

	if email != post.Email {
		err = errors.New("user email mismatch")
		log.Print(err)
		return err
	}

	err = s.commentRepo.Delete(id)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
