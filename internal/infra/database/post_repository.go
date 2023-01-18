package database

import (
	"gorm.io/gorm"
	"math"
	"nix_education/internal/domain"
)

const PostTableName = "posts"

type Post struct {
	Id     int64 `gorm:"primary_key;auto_increment;not_null"`
	UserId int64
	Title  string
	Body   string
}

type Posts struct {
	Items []Post
	Total uint64
	Pages uint64
}

//go:generate mockery --dir . --name PostRepository --output ./mocks
type PostRepository interface {
	Save(post Post) (Post, error)
	Find(id int64) (Post, error)
	FindAll(p domain.Pagination) (Posts, error)
	Update(post Post) (Post, error)
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

func (r postRepository) Save(p Post) (Post, error) {
	err := r.sess.Table(PostTableName).Create(&p).Error
	if err != nil {
		return Post{}, err
	}
	return p, nil
}

func (r postRepository) Find(id int64) (Post, error) {
	var p Post
	err := r.sess.Table(PostTableName).First(&p, "id = ?", id).Error
	if err != nil {
		return Post{}, err
	}
	return p, nil
}

func (r postRepository) FindAll(p domain.Pagination) (Posts, error) {
	var posts []Post
	offset := (p.Page - 1) * p.CountPerPage
	queryBuider := r.sess.Table(PostTableName).Limit(int(p.CountPerPage)).Offset(int(offset))
	err := queryBuider.Find(&posts).Error
	if err != nil {
		return Posts{}, err
	}
	pCol := mapToPostCollection(posts)
	result := r.sess.Table(PostTableName).Find(&posts)
	if result.Error != nil {
		return Posts{}, result.Error
	}
	total := result.RowsAffected
	pCol.Total = uint64(total)
	pCol.Pages = uint64(math.Ceil(float64(total) / float64(p.CountPerPage)))
	return pCol, nil
}

func (r postRepository) Update(p Post) (Post, error) {
	err := r.sess.Save(&p).Error
	if err != nil {
		return Post{}, err
	}
	return p, nil
}

func (r postRepository) Delete(id int64) error {
	err := r.sess.Table(PostTableName).Where("id = ?", id).Delete(Post{}).Error
	if err != nil {
		return err
	}
	return nil
}

func mapToPostCollection(posts []Post) Posts {
	res := Posts{
		Items: posts,
	}
	return res
}
