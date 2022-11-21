//go:build integration
// +build integration

package controllers_test

import (
	"NIX/internal/domain"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"net/http"
)

var userControllerTests = []*requestTest{
	{
		"Register",
		func(req *http.Request, migrator *migrate.Migrate) {
			resetDB(migrator)
		},
		"/api/v1/auth/register",
		"POST",
		`{"email":"email@example.com","password":"12345678","name":"User Name"}`,
		http.StatusCreated,
		`{"token":".{150,256}","user":{"id":1,"email":"email@example.com","name":"User Name"}}`,
		"wrong register new user response body",
	},
	{
		"Login",
		func(req *http.Request, migrator *migrate.Migrate) {},
		"/api/v1/auth/login",
		"POST",
		`{"email":"email@example.com","password":"12345678"}`,
		http.StatusOK,
		`{"token":".{150,256}","user":{"id":1,"email":"email@example.com","name":"User Name"}}`,
		"wrong login response body",
	},
}

func userModelMocker(n int) []domain.User {
	users := make([]domain.User, 0, n)
	for i := 1; i <= n; i++ {
		uModel := domain.User{
			Email:    fmt.Sprintf("email%d@example.com", i),
			Password: "123456",
			Name:     fmt.Sprintf("User Name %d", i),
		}

		user, err := userService.Save(uModel)
		if err != nil {
			log.Fatalf("userModelMocker() failed: %s", err)
		}
		users = append(users, user)
	}
	return users
}
