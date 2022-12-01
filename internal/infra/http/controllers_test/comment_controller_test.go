//go:build integration
// +build integration

package controllers_test_test

import (
	"NIX/internal/domain"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"net/http"
)

var commentControllerTests = []*requestTest{
	{
		"Get all comments for post",
		func(req *http.Request, migrator *migrate.Migrate) {
			resetDB(migrator)
			userModelMocker(2)
			postModelMocker(4, 1)
			commentModelMocker(4, 1, "email1@example.com")
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments",
		"GET",
		``,
		http.StatusOK,
		`{"items":\[(?:{"id":[1-4],"post_id":1,"name":"Comment[1-4] to post1","body":"Some test comment[1-4]"},?){4}\],"total":4,"pages":1}`,
		"wrong get all response body",
	},

	{
		"Find comment by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments/1",
		"GET",
		``,
		http.StatusOK,
		`{"id":1,"post_id":1,"name":"Comment1 to post1","body":"Some test comment1"}`,
		"wrong find comment response",
	},

	{
		"Save comment",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/2/comments",
		"POST",
		`{"name":"Comment1 to post2","body":"Some test comment5"}`,
		http.StatusCreated,
		`{"id":5,"post_id":2,"name":"Comment1 to post2","body":"Some test comment5"}`,
		"wrong save comment response body",
	},

	{
		"Update",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/2/comments/5",
		"PUT",
		`{"name":"Comment upd","body":"Some test upd"}`,
		http.StatusOK,
		`{"id":5,"post_id":2,"name":"Comment upd","body":"Some test upd"}`,
		"wrong update comment response",
	},

	{
		"Delete",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments/4",
		"DELETE",
		``,
		http.StatusOK,
		``,
		"wrong delete comment response",
	},

	{
		"Not found comment by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments/10",
		"GET",
		``,
		http.StatusNotFound,
		`{}`,
		"wrong not found comment response",
	},

	{
		"Find incorrect post comment by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/1/comments/fgfg",
		"GET",
		``,
		http.StatusBadRequest,
		`{"Func":"ParseInt","Num":"fgfg","Err":{}}`,
		"wrong find incorrect post comment by id response",
	},

	{
		"Save comment without require field (name)",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments",
		"POST",
		`{"body":"bad comment"}`,
		http.StatusUnprocessableEntity,
		`{}`,
		"wrong save comment without require field (name) response body",
	},

	{
		"Update another user comment",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/1/comments/1",
		"PUT",
		`{"title":"Post11","body":"test body"}`,
		http.StatusUnprocessableEntity,
		`{}`,
		"wrong update another user comment response",
	},

	{
		"Update not found comment by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments/10",
		"PUT",
		`{"title":"Post new","body":"test body"}`,
		http.StatusUnprocessableEntity,
		`{}`,
		"wrong update not found comment by id response",
	},

	{
		"Update incorrect type of id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments/dgdgdg",
		"PUT",
		`{"title":"Post new","body":"test body"}`,
		http.StatusBadRequest,
		`{"Func":"ParseInt","Num":"dgdgdg","Err":{}}`,
		"wrong update incorrect type of id response",
	},

	{
		"Update comment without require field (name)",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/1/comments/5",
		"PUT",
		`{"body":"test body"}`,
		http.StatusUnprocessableEntity,
		`{}`,
		"wrong update comment without require field (Title) response",
	},

	{
		"Delete another user comment",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 2, "email2@example.com")
		},
		"/api/v1/posts/1/comments/1",
		"DELETE",
		``,
		http.StatusNotFound,
		`{}`,
		"wrong delete another user comment response",
	},

	{
		"Delete not found comment by id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments/10",
		"DELETE",
		``,
		http.StatusNotFound,
		`{}`,
		"wrong delete not found comment by id response",
	},

	{
		"Delete incorrect type of id",
		func(req *http.Request, migrator *migrate.Migrate) {
			HeaderTokenMock(req, 1, "email1@example.com")
		},
		"/api/v1/posts/1/comments/sfsgs",
		"DELETE",
		``,
		http.StatusBadRequest,
		`{"Func":"ParseInt","Num":"sfsgs","Err":{}}`,
		"wrong delete incorrect type of id response",
	},
}

func commentModelMocker(n, id int, email string) {
	for i := 1; i <= n; i++ {

		postID := int64(id)

		cModel := domain.Comment{
			PostId: postID,
			Email:  email,
			Name:   fmt.Sprintf("Comment%d to post%d", i, postID),
			Body:   fmt.Sprintf("Some test comment%d", i),
		}

		_, err := commentService.Save(cModel)
		if err != nil {
			log.Fatalf("commentModelMocker() failed: %s", err)
		}
	}
}
