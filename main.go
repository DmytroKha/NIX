package main

import (
	"NIX/config"
	_ "NIX/docs"
	"NIX/internal/app"
	"NIX/internal/infra/database"
	"NIX/internal/infra/http/controllers"
	"NIX/internal/infra/http/router"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	mysqlG "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
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

	//BEGINNER. 1.	Налаштувати середовище розробки.
	//beginner.printHello()

	//BEGINNER. 2.	Робота з репозиторієм.
	//https://github.com/DmytroKha/NIX

	//BEGINNER. 3.	Отримання інформації з мережі. Є сервіс https://jsonplaceholder.typicode.com/ .
	//представляє REST API для отримання даних у форматі JSON. Сайт надає доступ до таких ресурсів:
	//beginner.getNetInformation()

	//BEGINNER. 4.	Горутини.
	//beginner.useGoroutine()

	//BEGINNER. 5.	Файлова система
	//beginner.useFileSystem()

	//BEGINNER. 6.	Робота с БД
	//beginner.useDB()

	//TRAINEE. 1.	Сodestyle
	//golangci-lint run

	//TRAINEE. 2.	Gitflow
	//???

	//TRAINEE. 3.	GORM
	//trainee.useDBWithGORM()

	//TRAINEE. 4.	Створення REST API
	//trainee.createRESTAPI()

	//TRAINEE. 5.	Echo framework
	echoRESTAPI()

	//TRAINEE. 6.	Swagger specification
	//Додай swagger до API. Використовуй пакет - swag
	//http://localhost:8080/swagger/index.html

	//TRAINEE. 7. OAuth 2.0 	Додай можливість реєстрації, авторизації користувачів,
	//використовуючи стандарт JWT; Додай авторизацію + реєстрацію через Google використовуючи протокол OAuth2.0.
	//Тільки авторизовані користувачі можуть писати пости та залишати коментарі,
	//кожний пост і комент прив'язаний до будь-якого користувача. Використовуй бібліотеку — golang/oauth2

	//TRAINEE. 8. Тестування	Напиши тести для свого API. Використовуй стандартну бібліотеку для тестування - testing

}

func echoRESTAPI() {

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
