package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

/*
type Message struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
*/

func main() {
	//1.	Налаштувати середовище розробки.
	//2.	Робота з репозиторієм.
	fmt.Println("Hello,NIX Education")

	//3.	Отримання інформації з мережі. Є сервіс https://jsonplaceholder.typicode.com/ .
	//представляє REST API для отримання даних у форматі JSON. Сайт надає доступ до таких ресурсів:
	/*
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

		var m []Message

		err = dec.Decode(&m)

		if err != nil {
			log.Fatalln(err)
		}

		data, err := json.MarshalIndent(m, "", "    ")

		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\n", data)
	*/

	//4.	Горутини.
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

			for true {
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

	//5.	Файлова система

}
