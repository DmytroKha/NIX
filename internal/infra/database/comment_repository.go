package database

import (
	"gorm.io/gorm"
	"math"
	"nix_education/internal/domain"
)

const CommentTableName = "comments"

type Comment struct {
	Id     int64 `gorm:"primary_key;auto_increment;not_null"`
	PostId int64
	Name   string
	Email  string
	Body   string
}

type Comments struct {
	Items []Comment
	Total uint64
	Pages uint64
}

//go:generate mockery --dir . --name CommentRepository --output ./mocks
type CommentRepository interface {
	Save(comment Comment) (Comment, error)
	Find(postId, id int64) (Comment, error)
	FindAll(postId int64, p domain.Pagination) (Comments, error)
	Update(comment Comment) (Comment, error)
	Delete(id int64) error
}

type commentRepository struct {
	sess *gorm.DB
}

func NewCommentRepository(dbSession *gorm.DB) CommentRepository {
	return &commentRepository{
		sess: dbSession,
	}
}

func (r commentRepository) Save(c Comment) (Comment, error) {
	err := r.sess.Table(CommentTableName).Create(&c).Error
	if err != nil {
		return Comment{}, err
	}
	return c, nil
}

func (r commentRepository) Find(postId, id int64) (Comment, error) {
	var c Comment
	err := r.sess.Table(CommentTableName).First(&c, "id = ? AND post_id = ?", id, postId).Error
	if err != nil {
		return Comment{}, err
	}
	return c, nil
}

func (r commentRepository) FindAll(postId int64, p domain.Pagination) (Comments, error) {
	var comments []Comment
	offset := (p.Page - 1) * p.CountPerPage
	queryBuider := r.sess.Table(CommentTableName).Limit(int(p.CountPerPage)).Offset(int(offset))
	err := queryBuider.Where("post_id = ?", postId).Find(&comments).Error
	if err != nil {
		return Comments{}, err
	}
	cCol := mapToCommentCollection(comments)
	result := r.sess.Table(CommentTableName).Where("post_id = ?", postId).Find(&comments)
	if result.Error != nil {
		return Comments{}, result.Error
	}
	total := result.RowsAffected
	cCol.Total = uint64(total)
	cCol.Pages = uint64(math.Ceil(float64(total) / float64(p.CountPerPage)))
	return cCol, nil
}

func (r commentRepository) Update(c Comment) (Comment, error) {
	err := r.sess.Save(&c).Error
	if err != nil {
		return Comment{}, err
	}
	return c, nil
}

func (r commentRepository) Delete(id int64) error {
	err := r.sess.Table(CommentTableName).Where("id = ?", id).Delete(Comment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func mapToCommentCollection(comments []Comment) Comments {
	res := Comments{
		Items: comments,
	}
	return res
}
