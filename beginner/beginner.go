package beginner

import (
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

//nolint
func printHello() {
	fmt.Println("Hello,NIX Education")
}

//nolint
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

//nolint
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

//nolint
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

//nolint
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
