package trainee_pt_one

import (
	"NIX/internal/domain"
	"fmt"
	"github.com/goccy/go-json"
	mysqlG "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"regexp"
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

type postHendler struct {
	store *gorm.DB
}

type commentHendler struct {
	store *gorm.DB
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

	var m []domain.Post

	err = msg.Decode(&m)

	if err != nil {
		log.Fatalln(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(m))

	for i := range m {
		ii := i
		go func(ii int) {

			db.Table("posts").Clauses(clause.OnConflict{DoNothing: true}).Create(&domain.Post{UserId: m[ii].UserId, Id: m[ii].Id, Title: m[ii].Title, Body: m[ii].Body})

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

			var c []domain.Comment

			err = com.Decode(&c)

			if err != nil {
				log.Fatalln(err)
			}
			wg.Add(len(c))

			for j := range c {
				jj := j
				go func(jj int) {
					db.Table("comments").Clauses(clause.OnConflict{DoNothing: true}).Create(&domain.Comment{PostId: c[jj].PostId, Id: c[jj].Id, Name: c[jj].Name, Email: c[jj].Email, Body: c[jj].Body})
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

	var posts []domain.Post
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

	var comments []domain.Comment
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

	var post domain.Post
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

	var comment domain.Comment
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

	var post, postNew domain.Post
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

	var comment, commentNew domain.Comment
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

	var post, postUpd domain.Post
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

	var comment, commentUpd domain.Comment
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

	var post domain.Post

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

	var comment domain.Comment

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
