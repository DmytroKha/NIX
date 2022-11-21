//go:build integration
// +build integration

package controllers_test

import (
	"github.com/golang-migrate/migrate/v4"
	"net/http"
)

var authControllerTests = []*requestTest{
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
