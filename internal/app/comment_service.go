package app

import (
	"NIX/internal/domain"
	"NIX/internal/infra/database"
	"log"
)

type CommentService interface {
	Save(comment domain.Comment) (domain.Comment, error)
	Update(comment domain.Comment) (domain.Comment, error)
	Find(id int64) (domain.Comment, error)
	FindAll(postId int64, p domain.Pagination) (domain.Comments, error)
	Delete(id int64) error
}

type commentService struct {
	commentRepo database.CommentRepository
}

func NewCommentService(r database.CommentRepository) CommentService {
	return commentService{
		commentRepo: r,
	}
}

func (s commentService) Save(p domain.Comment) (domain.Comment, error) {
	comment, err := s.commentRepo.Save(p)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}

	return comment, nil
}

func (s commentService) Find(id int64) (domain.Comment, error) {
	comment, err := s.commentRepo.Find(id)
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
	comment, err := s.commentRepo.Update(p)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}

	return comment, nil
}

func (s commentService) Delete(id int64) error {
	err := s.commentRepo.Delete(id)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
