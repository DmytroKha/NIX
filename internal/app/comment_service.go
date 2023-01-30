package app

import (
	"errors"
	"log"
	"nix_education/internal/domain"
	"nix_education/internal/infra/database"
	"nix_education/internal/infra/http/resources"
)

//go:generate mockery --dir . --name CommentService --output ./mocks
type CommentService interface {
	Save(comment database.Comment) (resources.CommentDto, error)
	Update(comment database.Comment) (resources.CommentDto, error)
	Find(postId, id int64) (resources.CommentDto, error)
	FindAll(postId int64, p domain.Pagination) (resources.CommentsDto, error)
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

func (s commentService) Save(p database.Comment) (resources.CommentDto, error) {
	_, err := s.postServise.Find(p.PostId)
	if err != nil {
		log.Print(err)
		return resources.CommentDto{}, err
	}
	comment, err := s.commentRepo.Save(p)
	if err != nil {
		log.Print(err)
		return resources.CommentDto{}, err
	}
	var commentDto resources.CommentDto
	return commentDto.DatabaseToDto(comment), nil
}

func (s commentService) Find(postId, id int64) (resources.CommentDto, error) {
	comment, err := s.commentRepo.Find(postId, id)
	if err != nil {
		log.Print(err)
		return resources.CommentDto{}, err
	}
	var commentDto resources.CommentDto
	return commentDto.DatabaseToDto(comment), nil
}

func (s commentService) FindAll(postId int64, p domain.Pagination) (resources.CommentsDto, error) {
	comments, err := s.commentRepo.FindAll(postId, p)
	if err != nil {
		log.Print(err)
		return resources.CommentsDto{}, err
	}
	var commentDto resources.CommentDto
	return commentDto.DatabaseToDtoCollection(comments), nil
}

func (s commentService) Update(p database.Comment) (resources.CommentDto, error) {
	post, err := s.commentRepo.Find(p.PostId, p.Id)
	if err != nil {
		log.Print(err)
		return resources.CommentDto{}, err
	}
	if p.Email != post.Email {
		err = errors.New("user email mismatch")
		return resources.CommentDto{}, err
	}
	comment, err := s.commentRepo.Update(p)
	if err != nil {
		log.Print(err)
		return resources.CommentDto{}, err
	}
	var commentDto resources.CommentDto
	return commentDto.DatabaseToDto(comment), nil
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
