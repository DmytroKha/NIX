//go:build integration
// +build integration

package controllers_test

import (
	"NIX/config"
	"NIX/internal/app"
	"NIX/internal/domain"
	"NIX/internal/infra/database"
	"NIX/internal/infra/http/controllers"
	"NIX/internal/infra/http/router"
	"bytes"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"
	mysqlG "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var authService app.AuthService
var userService app.UserService
var postService app.PostService
var commentService app.CommentService

type requestTest struct {
	name          string
	init          func(*http.Request, *migrate.Migrate) // Executed before test
	url           string
	method        string
	bodyData      string
	expectedCode  int
	responseRegex string
	msg           string // Test error message
}

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestControllers(t *testing.T) {

	var conf = config.Configuration{
		DatabaseName:        "nix_education",
		DatabaseHost:        "127.0.0.1:3306",
		DatabaseUser:        "root",
		DatabasePassword:    "root",
		MigrateToVersion:    "latest",
		MigrationLocation:   "../../database/migrations",
		FileStorageLocation: "",
		JwtSecret:           "1234567890",
		JwtTTL:              24 * time.Hour,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.DatabaseUser,
		conf.DatabasePassword,
		conf.DatabaseHost,
		conf.DatabaseName)
	db, err := gorm.Open(mysqlG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	connString := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s)/%s",
		conf.DatabaseUser,
		conf.DatabasePassword,
		conf.DatabaseHost,
		conf.DatabaseName,
	)

	migrator, err := migrate.New(
		"file://"+conf.MigrationLocation,
		connString)
	if err != nil {
		log.Fatalf("Unable to create Migrator: %q\n", err)
	}

	userRepository := database.NewUserRepository(db)
	userService = app.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	authService = app.NewAuthService(userService, conf)
	authController := controllers.NewAuthController(authService, userService)

	postRepository := database.NewPostRepository(db)
	postService = app.NewPostService(postRepository)
	postController := controllers.NewPostController(postService)

	commentRepository := database.NewCommentRepository(db)
	commentService = app.NewCommentService(commentRepository, postService)
	commentController := controllers.NewCommentController(commentService)

	// Create routes
	e := router.New(
		userController,
		authController,
		postController,
		commentController,
		conf)

	iterateOverTests(t, "AuthController", authControllerTests, e, migrator)
	//iterateOverTests(t, "UserController", userControllerTests, e, migrator)
	iterateOverTests(t, "PostController", postControllerTests, e, migrator)
	iterateOverTests(t, "CommentController", commentControllerTests, e, migrator)

}

func iterateOverTests(t *testing.T, name string, tests []*requestTest, router http.Handler, migrator *migrate.Migrate) {
	for _, tt := range tests {
		t.Run(name+" "+tt.name, func(t *testing.T) {
			fmt.Printf("[%-6s] %-35s", tt.method, tt.url)
			bodyData := tt.bodyData
			req, err := http.NewRequest(tt.method, tt.url, bytes.NewBufferString(bodyData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			tt.init(req, migrator)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			body := w.Body.String()

			fmt.Printf("[%d]\n", w.Code)
			require.Equal(t, tt.expectedCode, w.Code, "Response Status - "+tt.msg+"\nBody:\n"+body)
			require.Regexp(t, tt.responseRegex, body, "Response Content - "+tt.msg)
		})
	}
}

func HeaderTokenMock(req *http.Request, uId int64, email string) {
	tokenString, _ := authService.GenerateJwt(domain.User{Id: uId, Email: email})
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenString))
}

func resetDB(migrator *migrate.Migrate) {
	err := migrator.Down()
	if err != nil {
		log.Printf("migrator down: %q\n", err)
	}
	err = migrator.Up()
	if err != nil {
		log.Printf("migrator up: %q\n", err)
	}
}
