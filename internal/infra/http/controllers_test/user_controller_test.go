//go:build integration
// +build integration

package controllers_test_test

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"net/http"
	"nix_education/internal/infra/database"
)

var userControllerTests = []*requestTest{
	{
		"Set user password",
		func(req *http.Request, migrator *migrate.Migrate) {
			resetDB(migrator)
			userModelMocker(1, "")
		},
		"/api/v1/users/1",
		"PUT",
		`{"email":"email1@example.com","password":"12345678","name":"User Name 1"}`,
		http.StatusOK,
		`{"id":1,"email":"email1@example.com","name":"User Name 1"}`,
		"wrong set user password response body",
	},
}

func userModelMocker(n int, p string) []database.User {
	users := make([]database.User, 0, n)
	for i := 1; i <= n; i++ {
		uModel := database.User{
			Email:    fmt.Sprintf("email%d@example.com", i),
			Password: p,
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
