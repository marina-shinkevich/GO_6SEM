package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Метод: %s, Путь: %s", r.Method, r.URL.Path)

	fmt.Fprintf(w, "Запрос получен: %s %s", r.Method, r.URL.Path)
}

func main() {

	http.HandleFunc("/", handler)

	log.Println("Сервер GO03_01 запущен на порту :3000")
err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
