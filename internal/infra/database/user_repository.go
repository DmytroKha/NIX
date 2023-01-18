package database

import (
	"gorm.io/gorm"
)

const UserTableName = "users"

type User struct {
	Id       int64 `gorm:"primary_key;auto_increment;not_null"`
	Name     string
	Email    string
	Password string
}

type Users struct {
	Items []User
	Total uint64
	Pages uint64
}

//go:generate mockery --dir . --name UserRepository --output ./mocks
type UserRepository interface {
	Save(user User) (User, error)
	Update(user User) (User, error)
	FindByEmail(email string) (User, error)
}

type userRepository struct {
	sess *gorm.DB
}

func NewUserRepository(dbSession *gorm.DB) UserRepository {
	return &userRepository{
		sess: dbSession,
	}
}

func (r userRepository) Save(u User) (User, error) {
	err := r.sess.Table(UserTableName).Create(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r userRepository) Update(u User) (User, error) {
	err := r.sess.Save(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *userRepository) FindByEmail(email string) (User, error) {
	var u User
	err := r.sess.Table(UserTableName).First(&u, "email = ?", email).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}
