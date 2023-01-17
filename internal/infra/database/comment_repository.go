package database

import (
	"gorm.io/gorm"
	"math"
	"nix_education/internal/domain"
)

const CommentTableName = "comments"

type comment struct {
	Id     int64 `gorm:"primary_key;auto_increment;not_null"`
	PostId int64
	Name   string
	Email  string
	Body   string
}

//go:generate mockery --dir . --name CommentRepository --output ./mocks
type CommentRepository interface {
	Save(comment domain.Comment) (domain.Comment, error)
	Find(postId, id int64) (domain.Comment, error)
	FindAll(postId int64, p domain.Pagination) (domain.Comments, error)
	Update(comment domain.Comment) (domain.Comment, error)
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

func (r commentRepository) Save(p domain.Comment) (domain.Comment, error) {
	var cmt comment

	cmt.FromDomainModel(p)
	err := r.sess.Table(CommentTableName).Create(&cmt).Error
	if err != nil {
		return domain.Comment{}, err
	}

	return cmt.ToDomainModel(), nil
}

func (r commentRepository) Find(postId, id int64) (domain.Comment, error) {
	var cmt comment

	err := r.sess.Table(CommentTableName).First(&cmt, "id = ? AND post_id = ?", id, postId).Error
	if err != nil {
		return domain.Comment{}, err
	}

	return cmt.ToDomainModel(), nil
}

func (r commentRepository) FindAll(postId int64, p domain.Pagination) (domain.Comments, error) {
	var comments []comment
	offset := (p.Page - 1) * p.CountPerPage
	queryBuider := r.sess.Table(CommentTableName).Limit(int(p.CountPerPage)).Offset(int(offset))
	err := queryBuider.Where("post_id = ?", postId).Find(&comments).Error
	if err != nil {
		return domain.Comments{}, err
	}

	dcomments := mapToCommentDomainCollection(comments)

	result := r.sess.Table(CommentTableName).Where("post_id = ?", postId).Find(&comments)
	if result.Error != nil {
		return domain.Comments{}, result.Error
	}

	total := result.RowsAffected

	dcomments.Total = uint64(total)
	dcomments.Pages = uint64(math.Ceil(float64(total) / float64(p.CountPerPage)))

	return dcomments, nil
}

func (r commentRepository) Update(p domain.Comment) (domain.Comment, error) {
	var cmt comment

	cmt.FromDomainModel(p)

	err := r.sess.Save(&cmt).Error
	if err != nil {
		return domain.Comment{}, err
	}

	return cmt.ToDomainModel(), nil
}

func (r commentRepository) Delete(id int64) error {
	err := r.sess.Table(CommentTableName).Where("id = ?", id).Delete(domain.Comment{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (p comment) ToDomainModel() domain.Comment {
	return domain.Comment{
		Id:     p.Id,
		PostId: p.PostId,
		Name:   p.Name,
		Email:  p.Email,
		Body:   p.Body,
	}
}

func (p *comment) FromDomainModel(dp domain.Comment) {
	p.Id = dp.Id
	p.PostId = dp.PostId
	p.Name = dp.Name
	p.Email = dp.Email
	p.Body = dp.Body
}

func mapToCommentDomainCollection(comments []comment) domain.Comments {
	var result []domain.Comment

	if len(comments) == 0 {
		result = make([]domain.Comment, 0)
	}

	for _, c := range comments {
		d := c.ToDomainModel()
		result = append(result, d)
	}

	res := domain.Comments{
		Items: result,
	}

	return res
}
