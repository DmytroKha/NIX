package app

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"nix_education/config"
	"nix_education/internal/infra/database"
	"nix_education/internal/infra/http/resources"
	"strconv"
	"time"
)

//go:generate mockery --dir . --name AuthService --output ./mocks
type AuthService interface {
	Register(user database.User) (database.User, string, error)
	Login(user database.User) (database.User, string, error)
	LoginGoogle(email string) (database.User, string, error)
	GenerateJwt(user database.User) (string, error)
}

type authService struct {
	userService UserService
	config      config.Configuration
}

func NewAuthService(us UserService, cf config.Configuration) AuthService {
	return authService{
		userService: us,
		config:      cf,
	}
}

func (s authService) Register(u database.User) (database.User, string, error) {
	_, err := s.userService.FindByEmail(u.Email)
	if err == nil {
		log.Printf("invalid credentials")
		return database.User{}, "", errors.New("invalid credentials")
	}
	user, err := s.userService.Save(u)
	if err != nil {
		log.Print(err)
		return database.User{}, "", err
	}
	token, err := s.GenerateJwt(user)
	return user, token, err
}

func (s authService) Login(user database.User) (database.User, string, error) {
	u, err := s.userService.FindByEmail(user.Email)
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	valid := s.checkPasswordHash(user.Password, u.Password)
	if !valid {
		return database.User{}, "", errors.New("invalid credentials")
	}
	token, err := s.GenerateJwt(u)
	return u, token, err
}

func (s authService) LoginGoogle(email string) (database.User, string, error) {
	u, err := s.userService.FindByEmail(email)
	if err != nil {
		u.Email = email
		u, err = s.userService.Save(u)
		if err != nil {
			log.Print(err)
			return database.User{}, "", err
		}
	}
	token, err := s.GenerateJwt(u)
	return u, token, err
}

func (s authService) GenerateJwt(user database.User) (string, error) {
	claims := resources.JwtClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(int(user.Id)),
			ExpiresAt: time.Now().Add(s.config.JwtTTL).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte(s.config.JwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s authService) checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
