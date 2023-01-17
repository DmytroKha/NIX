package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	mysqlG "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"nix_education/config"
	_ "nix_education/docs"
	"nix_education/internal/app"
	"nix_education/internal/infra/database"
	"nix_education/internal/infra/http/controllers"
	"nix_education/internal/infra/http/router"
)

// @title       NIX_Education API
// @version     1.0
// @description API Server for NIX_Education application.

// @host     localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
func main() {
	var conf = config.GetConfiguration()

	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v",
		conf.DatabaseUser,
		conf.DatabasePassword,
		conf.DatabaseHost,
		conf.DatabaseName)
	db, err := gorm.Open(mysqlG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	userRepository := database.NewUserRepository(db)
	userService := app.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	authService := app.NewAuthService(userService, conf)
	authController := controllers.NewAuthController(authService, userService)

	postRepository := database.NewPostRepository(db)
	postService := app.NewPostService(postRepository)
	postController := controllers.NewPostController(postService)

	commentRepository := database.NewCommentRepository(db)
	commentService := app.NewCommentService(commentRepository, postService)
	commentController := controllers.NewCommentController(commentService)

	e := router.New(
		userController,
		authController,
		postController,
		commentController,
		conf)

	// service start at port :8080
	err = e.Start(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}
