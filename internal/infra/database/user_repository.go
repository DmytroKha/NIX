package database

import (
	"gorm.io/gorm"
	"nix_education/internal/domain"
)

const UserTableName = "users"

type user struct {
	Id       int64 `gorm:"primary_key;auto_increment;not_null"`
	Name     string
	Email    string
	Password string
}

type UserRepository interface {
	Save(user domain.User) (domain.User, error)
	Update(user domain.User) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
}

type userRepository struct {
	sess *gorm.DB
}

func NewUserRepository(dbSession *gorm.DB) UserRepository {
	return &userRepository{
		sess: dbSession,
	}
}

func (r userRepository) Save(u domain.User) (domain.User, error) {
	var usr user

	usr.FromDomainModel(u)

	err := r.sess.Table(UserTableName).Create(&usr).Error
	if err != nil {
		return domain.User{}, err
	}

	return usr.ToDomainModel(), nil
}

func (r userRepository) Update(u domain.User) (domain.User, error) {
	var usr user

	usr.FromDomainModel(u)

	err := r.sess.Save(&usr).Error
	if err != nil {
		return domain.User{}, err
	}

	return usr.ToDomainModel(), nil
}

func (r *userRepository) FindByEmail(email string) (domain.User, error) {
	var usr user

	err := r.sess.Table(UserTableName).First(&usr, "email = ?", email).Error
	if err != nil {
		return domain.User{}, err
	}

	return usr.ToDomainModel(), nil
}

func (u user) ToDomainModel() domain.User {
	return domain.User{
		Id:       u.Id,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *user) FromDomainModel(du domain.User) {
	u.Id = du.Id
	u.Name = du.Name
	u.Email = du.Email
	u.Password = du.Password
}
