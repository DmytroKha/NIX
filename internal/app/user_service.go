package app

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"nix_education/internal/infra/database"
	"nix_education/internal/infra/http/requests"
)

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Save(user database.User) (database.User, error)
	SetPassword(usr requests.UserRequest) (database.User, error)
	GeneratePasswordHash(password string) (string, error)
	FindByEmail(email string) (database.User, error)
}

type userService struct {
	userRepo database.UserRepository
}

func NewUserService(r database.UserRepository) UserService {
	return userService{
		userRepo: r,
	}
}

func (s userService) Save(u database.User) (database.User, error) {
	var err error
	u.Password, err = s.GeneratePasswordHash(u.Password)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	user, err := s.userRepo.Save(u)
	if err != nil {
		log.Print(err)
		return database.User{}, err
	}
	return user, nil
}

func (s userService) SetPassword(usr requests.UserRequest) (database.User, error) {
	u, err := usr.ToDatabaseModel()
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}

	user, err := s.FindByEmail(u.Email)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(""))
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	user.Password, err = s.GeneratePasswordHash(u.Password)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Print(err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) FindByEmail(email string) (database.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	return user, err
}

func (s userService) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
