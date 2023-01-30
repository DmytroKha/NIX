//go:build integration
// +build integration

package controllers_test_test

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"net/http"
	"nix_education/internal/infra/http/requests"
)

var postControllerTests = []*requestTest{
	{
		"Get all posts for user",
		func(req *http.Request, migrator *migrate.Migrate) {
			resetDB(migrator)
			userModelMocker(2, "123456")
			postModelMocker(4, 1)
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts",
		"GET",
		``,
		http.StatusOK,
		`{"items":\[(?:{"id":[1-4],"user_id":1,"title":"Post[1-4]","body":"some post text[1-4]"},?){4}\],"total":4,"pages":1}`,
		"wrong get all response body",
	},
	{
		"Find post by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/2",
		"GET",
		``,
		http.StatusOK,
		`{"id":2,"user_id":1,"title":"Post2","body":"some post text2"}`,
		"wrong find post response",
	},

	{
		"Save post",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts",
		"POST",
		`{"title":"Post5","body":"some post text5"}`,
		http.StatusCreated,
		`{"id":5,"user_id":2,"title":"Post5","body":"some post text5"}`,
		"wrong save post response body",
	},

	{
		"Update",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/4",
		"PUT",
		`{"title":"Post6","body":"post 4 to post 6"}`,
		http.StatusOK,
		`{"id":4,"user_id":1,"title":"Post6","body":"post 4 to post 6"}`,
		"wrong update post response",
	},
	{
		"Delete",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/4",
		"DELETE",
		``,
		http.StatusOK,
		``,
		"wrong delete post response",
	},

	{
		"Not found post by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/10",
		"GET",
		``,
		http.StatusNotFound,
		`{}`,
		"wrong not found post response",
	},

	{
		"Find incorrect user post by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/fgfg",
		"GET",
		``,
		http.StatusBadRequest,
		`{"Func":"ParseInt","Num":"fgfg","Err":{}}`,
		"wrong find incorrect user post by id response",
	},

	{
		"Save pet without require field (title)",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts",
		"POST",
		`{"body":"bad post"}`,
		http.StatusUnprocessableEntity,
		`{}`,
		"wrong save post without require field (title) response body",
	},

	{
		"Update another user pet",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/1",
		"PUT",
		`{"title":"Post11","body":"test body"}`,
		http.StatusInternalServerError,
		`{}`,
		"wrong update another user post response",
	},

	{
		"Update not found pet by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/10",
		"PUT",
		`{"title":"Post new","body":"test body"}`,
		http.StatusInternalServerError,
		`{}`,
		"wrong update not found post by id response",
	},

	{
		"Update incorrect type of id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/dgdgdg",
		"PUT",
		`{"title":"Post new","body":"test body"}`,
		http.StatusBadRequest,
		`{"Func":"ParseInt","Num":"dgdgdg","Err":{}}`,
		"wrong update incorrect type of id response",
	},

	{
		"Update post without require field (tite)",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/5",
		"PUT",
		`{"body":"test body"}`,
		http.StatusUnprocessableEntity,
		`{}`,
		"wrong update post without require field (Title) response",
	},

	{
		"Delete another user post",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/1",
		"DELETE",
		``,
		http.StatusNotFound,
		`{}`,
		"wrong delete another user post response",
	},

	{
		"Delete not found post by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/10",
		"DELETE",
		``,
		http.StatusNotFound,
		`{}`,
		"wrong delete not found post by id response",
	},

	{
		"Delete incorrect type of id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/pets/sfsgs",
		"DELETE",
		``,
		http.StatusNotFound,
		`{"message":"Not Found"}`,
		"wrong delete incorrect type of id response",
	},
}

func postModelMocker(n, id int) {
	for i := 1; i <= n; i++ {
		userID := int64(id)
		/*
			pModel := database.Post{
				Title:  fmt.Sprintf("Post%d", i),
				Body:   fmt.Sprintf("some post text%d", i),
				UserId: userID,
			}
		*/
		pModel := requests.PostRequest{
			Title: fmt.Sprintf("Post%d", i),
			Body:  fmt.Sprintf("some post text%d", i),
		}
		_, err := postService.Save(pModel, userID)
		if err != nil {
			log.Fatalf("postModelMocker() failed: %s", err)
		}
	}
}
