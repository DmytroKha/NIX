package app

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"nix_education/internal/domain"
	"nix_education/internal/infra/database"
)

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Save(user domain.User) (domain.User, error)
	SetPassword(user domain.User) (domain.User, error)
	GeneratePasswordHash(password string) (string, error)
	FindByEmail(email string) (domain.User, error)
}

type userService struct {
	userRepo database.UserRepository
}

func NewUserService(r database.UserRepository) UserService {
	return userService{
		userRepo: r,
	}
}

func (s userService) Save(u domain.User) (domain.User, error) {
	var err error

	u.Password, err = s.GeneratePasswordHash(u.Password)
	if err != nil {
		log.Printf("UserService: %s", err)
		return domain.User{}, err
	}

	user, err := s.userRepo.Save(u)
	if err != nil {
		log.Print(err)
		return domain.User{}, err
	}

	return user, nil
}

func (s userService) SetPassword(u domain.User) (domain.User, error) {

	user, err := s.FindByEmail(u.Email)
	if err != nil {
		log.Printf("UserService: %s", err)
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(""))
	if err != nil {
		log.Printf("UserService: %s", err)
		return domain.User{}, err
	}

	user.Password, err = s.GeneratePasswordHash(u.Password)
	if err != nil {
		log.Printf("UserService: %s", err)
		return domain.User{}, err
	}

	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Print(err)
		return domain.User{}, err
	}

	return updatedUser, nil
}

func (s userService) FindByEmail(email string) (domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("UserService: %s", err)
		return domain.User{}, err
	}

	return user, err
}

func (s userService) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
