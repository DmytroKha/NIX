package database

import (
	"NIX/internal/domain"
	"gorm.io/gorm"
	"math"
)

const PostTableName = "posts"

type post struct {
	Id     int64 `gorm:"primary_key;auto_increment;not_null"`
	UserId int64
	Title  string
	Body   string
}

type PostRepository interface {
	Save(post domain.Post) (domain.Post, error)
	Find(id int64) (domain.Post, error)
	FindAll(p domain.Pagination) (domain.Posts, error)
	Update(post domain.Post) (domain.Post, error)
	Delete(id int64) error
}

type postRepository struct {
	sess *gorm.DB
}

func NewPostRepository(dbSession *gorm.DB) PostRepository {
	return &postRepository{
		sess: dbSession,
	}
}

func (r postRepository) Save(p domain.Post) (domain.Post, error) {
	var pst post

	pst.FromDomainModel(p)

	err := r.sess.Table(PostTableName).Create(&pst).Error
	if err != nil {
		return domain.Post{}, err
	}

	return pst.ToDomainModel(), nil
}

func (r postRepository) Find(id int64) (domain.Post, error) {
	var pst post

	err := r.sess.Table(PostTableName).First(&pst, "id = ?", id).Error
	if err != nil {
		return domain.Post{}, err
	}

	return pst.ToDomainModel(), nil
}

func (r postRepository) FindAll(p domain.Pagination) (domain.Posts, error) {
	var posts []post
	offset := (p.Page - 1) * p.CountPerPage
	queryBuider := r.sess.Table(PostTableName).Limit(int(p.CountPerPage)).Offset(int(offset))
	err := queryBuider.Find(&posts).Error
	if err != nil {
		return domain.Posts{}, err
	}

	dposts := mapToPostDomainCollection(posts)

	result := r.sess.Table(PostTableName).Find(&posts)
	if result.Error != nil {
		return domain.Posts{}, result.Error
	}

	total := result.RowsAffected

	dposts.Total = uint64(total)
	dposts.Pages = uint64(math.Ceil(float64(total) / float64(p.CountPerPage)))

	return dposts, nil
}

func (r postRepository) Update(p domain.Post) (domain.Post, error) {
	var pst post

	pst.FromDomainModel(p)

	err := r.sess.Save(&pst).Error
	if err != nil {
		return domain.Post{}, err
	}

	return pst.ToDomainModel(), nil
}

func (r postRepository) Delete(id int64) error {
	err := r.sess.Table(PostTableName).Where("id = ?", id).Delete(domain.Post{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (p post) ToDomainModel() domain.Post {
	return domain.Post{
		Id:     p.Id,
		UserId: p.UserId,
		Title:  p.Title,
		Body:   p.Body,
	}
}

func (p *post) FromDomainModel(dp domain.Post) {
	p.Id = dp.Id
	p.UserId = dp.UserId
	p.Title = dp.Title
	p.Body = dp.Body
}

func mapToPostDomainCollection(posts []post) domain.Posts {
	var result []domain.Post

	if len(posts) == 0 {
		result = make([]domain.Post, 0)
	}

	for _, c := range posts {
		d := c.ToDomainModel()
		result = append(result, d)
	}

	res := domain.Posts{
		Items: result,
	}

	return res
}
