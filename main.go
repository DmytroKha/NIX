package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Message []struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	//1.	Налаштувати середовище розробки.
	//2.	Робота з репозиторієм.
	fmt.Println("Hello,NIX Education")

	//3.	Отримання інформації з мережі. Є сервіс https://jsonplaceholder.typicode.com/ .
	//представляє REST API для отримання даних у форматі JSON. Сайт надає доступ до таких ресурсів:
	formData := strings.NewReader(`{}`)

	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", formData)

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

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("bad request")
		log.Fatalln(err)
	}

	dec := json.NewDecoder(resp.Body)

	var m Message
	err = dec.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}
	data2, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data2)

}
