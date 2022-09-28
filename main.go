package main

import (
	"NIX/internal/app"
	"NIX/internal/infra/database"
	"NIX/internal/infra/http/controllers"
	"NIX/internal/infra/http/requests"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	mysqlG "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	listPostRe   = regexp.MustCompile(`^\/posts[\/]*$`)
	getPostRe    = regexp.MustCompile(`^\/posts\/(\d+)$`)
	createPostRe = regexp.MustCompile(`^\/posts[\/]*$`)
	updatePostRe = regexp.MustCompile(`^\/posts\/(\d+)$`)
	deletePostRe = regexp.MustCompile(`^\/posts\/(\d+)$`)

	listCommentRe   = regexp.MustCompile(`^\/comments[\/]*$`)
	getCommentRe    = regexp.MustCompile(`^\/comments\/(\d+)$`)
	createCommentRe = regexp.MustCompile(`^\/comments[\/]*$`)
	updateCommentRe = regexp.MustCompile(`^\/comments\/(\d+)$`)
	deleteCommentRe = regexp.MustCompile(`^\/comments\/(\d+)$`)
)

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

type postHendler struct {
	store *gorm.DB
}

type commentHendler struct {
	store *gorm.DB
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {

	//BEGINNER. 1.	Налаштувати середовище розробки.
	//printHello()

	//BEGINNER. 2.	Робота з репозиторієм.
	//https://github.com/DmytroKha/NIX

	//BEGINNER. 3.	Отримання інформації з мережі. Є сервіс https://jsonplaceholder.typicode.com/ .
	//представляє REST API для отримання даних у форматі JSON. Сайт надає доступ до таких ресурсів:
	//getNetInformation()

	//BEGINNER. 4.	Горутини.
	//useGoroutine()

	//BEGINNER. 5.	Файлова система
	//useFileSystem()

	//BEGINNER. 6.	Робота с БД
	//useDB()

	//TRAINEE. 1.	Сodestyle
	//golangci-lint run

	//TRAINEE. 2.	Gitflow
	//???

	//TRAINEE. 3.	GORM
	//useDBWithGORM()

	//TRAINEE. 4.	Створення REST API
	//createRESTAPI()

	//TRAINEE. 5.	Echo framework
	echoRESTAPI()

	//TRAINEE. 6.	Swagger specification
	//Додай swagger до API. Використовуй пакет - swag

}

func printHello() {
	fmt.Println("Hello,NIX Education")
}

func getNetInformation() {
	tempData := strings.NewReader(`{}`)

	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", tempData)

	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("bad request")
		log.Fatalln(err)
	}

	dec := json.NewDecoder(resp.Body)

	var m []Post

	err = dec.Decode(&m)

	if err != nil {
		log.Fatalln(err)
	}

	data, err := json.MarshalIndent(m, "", "    ")

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", data)
}

func useGoroutine() {
	wg := new(sync.WaitGroup)
	var totalString []string
	n := 5
	wg.Add(n)

	for i := 1; i <= n; i++ {
		b := i
		go func(b int) {
			s := "https://jsonplaceholder.typicode.com/posts/" + strconv.Itoa(b)
			resp, err := http.Get(s)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			for {
				bs := make([]byte, 1014)
				j, err := resp.Body.Read(bs)
				totalString = append(totalString, string(bs[:j]))
				if j == 0 || err != nil {
					break
				}
			}
			defer wg.Done()
		}(b)
	}

	wg.Wait()
	for i := range totalString {
		fmt.Println(totalString[i])
	}
}

func useFileSystem() {
	wg := new(sync.WaitGroup)
	n := 5
	wg.Add(n)

	err := os.MkdirAll("storage/posts", 0755)
	if err != nil {
		// print it out
		fmt.Println(err)
	}

	for i := 1; i <= n; i++ {
		b := i
		go func(b int) {
			s := "https://jsonplaceholder.typicode.com/posts/" + strconv.Itoa(b)
			resp, err := http.Get(s)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			for {
				bs := make([]byte, 1014)
				j, err1 := resp.Body.Read(bs)

				mydata := bs[:j]
				err2 := ioutil.WriteFile("storage/posts/"+strconv.Itoa(b)+".txt", mydata, 0644)
				// Обработка ошибки
				if err2 != nil {
					// print it out
					fmt.Println(err)
				}

				if j == 0 || err1 != nil {
					break
				}

			}
			defer wg.Done()
		}(b)
	}

	wg.Wait()

}

func useDB() {

	s := "https://jsonplaceholder.typicode.com/posts?userId=7"
	resp, err := http.Get(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	msg := json.NewDecoder(resp.Body)

	var m []Post

	err = msg.Decode(&m)

	if err != nil {
		log.Fatalln(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(m))

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/nix_education")
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		_ = db.Close()
	}()

	for i := range m {
		ii := i
		go func(ii int) {

			//insertPost, errPost := db.Query("INSERT INTO `posts` (`userId`, `id`, `title`, `body`) VALUES (?,?,?,?)", m[ii].UserId, m[ii].Id, m[ii].Title, m[ii].Body)
			insertPost, errPost := db.Query("INSERT INTO posts (user_id, id, title, body) SELECT * FROM (SELECT ? AS user_id, ? AS id, ? AS title, ? AS body) AS new_value WHERE NOT EXISTS (SELECT id FROM posts WHERE id = ?) LIMIT 1", m[ii].UserId, m[ii].Id, m[ii].Title, m[ii].Body, m[ii].Id)
			if errPost != nil {
				log.Fatalln(errPost)
			}

			defer func() {
				_ = insertPost.Close()
			}()

			sс := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%v", m[ii].Id)
			respCom, errCom := http.Get(sс)
			if errCom != nil {
				fmt.Println(errCom)
				return
			}
			defer func() {
				_ = respCom.Body.Close()
			}()

			com := json.NewDecoder(respCom.Body)

			var c []Comment

			err = com.Decode(&c)

			if err != nil {
				log.Fatalln(err)
			}
			wg.Add(len(c))

			for j := range c {
				jj := j
				go func(jj int) {
					//insertComment, errNewCom := db.Query("INSERT INTO `comments` (`postId`, `id`, `name`, `email`, `body`) VALUES (?,?,?,?,?)", c[jj].PostId, c[jj].Id, c[jj].Name, c[jj].Email, c[jj].Body)
					insertComment, errNewCom := db.Query("INSERT INTO `comments` (`post_id`, `id`, `name`, `email`, `body`) SELECT * FROM (SELECT ? AS post_id, ? AS id, ? AS name, ? AS email, ? AS body) AS new_value WHERE NOT EXISTS (SELECT id FROM comments WHERE id = ?) LIMIT 1", c[jj].PostId, c[jj].Id, c[jj].Name, c[jj].Email, c[jj].Body, c[jj].Id)
					if errNewCom != nil {
						log.Fatalln(errNewCom)
					}
					defer func() {
						_ = insertComment.Close()
					}()
					defer wg.Done()
				}(jj)
			}

			defer wg.Done()
		}(ii)
	}

	wg.Wait()

}

func useDBWithGORM() {

	dsn := "root:root@tcp(127.0.0.1:3306)/nix_education"
	db, err := gorm.Open(mysqlG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	s := "https://jsonplaceholder.typicode.com/posts?userId=7"
	resp, err := http.Get(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	msg := json.NewDecoder(resp.Body)

	var m []Post

	err = msg.Decode(&m)

	if err != nil {
		log.Fatalln(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(m))

	for i := range m {
		ii := i
		go func(ii int) {

			db.Table("posts").Clauses(clause.OnConflict{DoNothing: true}).Create(&Post{UserId: m[ii].UserId, Id: m[ii].Id, Title: m[ii].Title, Body: m[ii].Body})

			sс := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?postId=%v", m[ii].Id)
			respCom, errCom := http.Get(sс)
			if errCom != nil {
				fmt.Println(errCom)
				return
			}
			defer func() {
				_ = respCom.Body.Close()
			}()

			com := json.NewDecoder(respCom.Body)

			var c []Comment

			err = com.Decode(&c)

			if err != nil {
				log.Fatalln(err)
			}
			wg.Add(len(c))

			for j := range c {
				jj := j
				go func(jj int) {
					db.Table("comments").Clauses(clause.OnConflict{DoNothing: true}).Create(&Comment{PostId: c[jj].PostId, Id: c[jj].Id, Name: c[jj].Name, Email: c[jj].Email, Body: c[jj].Body})
					defer wg.Done()
				}(jj)
			}
			defer wg.Done()
		}(ii)
	}

	wg.Wait()

}

func createRESTAPI() {
	dsn := "root:root@tcp(127.0.0.1:3306)/nix_education"
	db, err := gorm.Open(mysqlG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	postH := &postHendler{store: db.Table("posts")}
	commentH := &commentHendler{store: db.Table("comments")}

	mux.Handle("/posts", postH)
	mux.Handle("/posts/", postH)
	mux.Handle("/comments", commentH)
	mux.Handle("/comments/", commentH)

	http.ListenAndServe("localhost:8080", mux)
}

func (h *postHendler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listPostRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getPostRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPost && createPostRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodPut && updatePostRe.MatchString(r.URL.Path):
		h.Update(w, r)
		return
	case r.Method == http.MethodDelete && deletePostRe.MatchString(r.URL.Path):
		h.Delete(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func (h *commentHendler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listCommentRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getCommentRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPost && createCommentRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodPut && updateCommentRe.MatchString(r.URL.Path):
		h.Update(w, r)
		return
	case r.Method == http.MethodDelete && deleteCommentRe.MatchString(r.URL.Path):
		h.Delete(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func (h *postHendler) List(w http.ResponseWriter, r *http.Request) {

	var posts []Post
	result := h.store.Find(&posts)

	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("posts not found"))
		return
	}

	jsonBytes, err := json.Marshal(posts)
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *commentHendler) List(w http.ResponseWriter, r *http.Request) {

	var comments []Comment
	result := h.store.Find(&comments)

	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("comments not found"))
		return
	}

	jsonBytes, err := json.Marshal(comments)
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *postHendler) Get(w http.ResponseWriter, r *http.Request) {

	matches := getPostRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	var post Post
	h.store.First(&post, matches[1])

	if post.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("post not found"))
		return
	}

	jsonBytes, err := json.Marshal(post)
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *commentHendler) Get(w http.ResponseWriter, r *http.Request) {

	matches := getCommentRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	var comment Comment
	h.store.First(&comment, matches[1])

	if comment.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("comment not found"))
		return
	}

	jsonBytes, err := json.Marshal(comment)
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *postHendler) Create(w http.ResponseWriter, r *http.Request) {

	var post, postNew Post
	if err := json.NewDecoder(r.Body).Decode(&postNew); err != nil {
		internalServerError(w, r)
		return
	}

	h.store.First(&post, postNew.Id)

	if post.Id != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("post already exist"))
		return
	}

	result := h.store.Create(postNew)

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("post not created"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *commentHendler) Create(w http.ResponseWriter, r *http.Request) {

	var comment, commentNew Comment
	if err := json.NewDecoder(r.Body).Decode(&commentNew); err != nil {
		internalServerError(w, r)
		return
	}

	h.store.First(&comment, commentNew.Id)

	if comment.Id != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("comment already exist"))
		return
	}

	result := h.store.Create(commentNew)

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("comment not created"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *postHendler) Update(w http.ResponseWriter, r *http.Request) {

	matches := updatePostRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	var post, postUpd Post
	h.store.First(&post, matches[1])

	if post.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("post not found"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&postUpd); err != nil {
		internalServerError(w, r)
		return
	}

	postUpd.Id = post.Id
	postUpd.UserId = post.UserId

	h.store.Model(&post).Updates(&postUpd)

	jsonBytes, err := json.Marshal(post)
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *commentHendler) Update(w http.ResponseWriter, r *http.Request) {

	matches := updateCommentRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	var comment, commentUpd Comment
	h.store.First(&comment, matches[1])

	if comment.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("post not found"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&commentUpd); err != nil {
		internalServerError(w, r)
		return
	}

	commentUpd.Id = comment.Id
	commentUpd.PostId = comment.PostId

	h.store.Model(&comment).Updates(&commentUpd)

	jsonBytes, err := json.Marshal(comment)
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *postHendler) Delete(w http.ResponseWriter, r *http.Request) {

	matches := deletePostRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	var post Post

	result := h.store.Delete(&post, matches[1])

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("post not deleted"))
		return
	}

	if result.Error == nil && result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("post not found"))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *commentHendler) Delete(w http.ResponseWriter, r *http.Request) {

	matches := deleteCommentRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	var comment Comment

	result := h.store.Delete(&comment, matches[1])

	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("comment not deleted"))
		return
	}

	if result.Error == nil && result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("comment not found"))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func echoRESTAPI() {

	dsn := "root:root@tcp(127.0.0.1:3306)/nix_education"
	db, err := gorm.Open(mysqlG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	postRepository := database.NewPostRepository(db)
	postService := app.NewPostService(postRepository)
	postController := controllers.NewPostController(postService)

	commentRepository := database.NewCommentRepository(db)
	commentService := app.NewCommentService(commentRepository)
	commentController := controllers.NewCommentController(commentService)

	e := echo.New()
	e.Validator = requests.NewValidator()

	api := e.Group("/api/v1", serverHeader)
	api.GET("/posts", postController.FindAll)
	api.POST("/posts", postController.Save)
	api.GET("/posts/:id", postController.Find)
	api.PUT("/posts/:id", postController.Update) //розібратися з контекстом
	api.DELETE("/posts/:id", postController.Delete)

	api.GET("/comments", commentController.FindAll)
	api.POST("/comments", commentController.Save)
	api.GET("/comments/:id", commentController.Find)
	api.PUT("/comments/:id", commentController.Update)
	api.DELETE("/comments/:id", commentController.Delete)

	// service start at port :8080
	err = e.Start(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("x-version", "Test/v1.0")
		return next(c)
	}
}
