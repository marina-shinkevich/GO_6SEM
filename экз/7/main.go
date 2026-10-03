package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// handlerQ_XY — сработает, если в URL переданы ОБА параметра (и x, и y).
func handlerQ_XY(w http.ResponseWriter, r *http.Request) {
	// Извлекаем значения параметров из разобранного URL
	x := r.URL.Query().Get("x")
	y := r.URL.Query().Get("y")

	fmt.Fprintf(w, "Сработал маршрут XY!\nПараметры: x = %s, y = %s\n", x, y)
}

// handlerQ_X — сработает, если в URL передан ТОЛЬКО параметр x.
func handlerQ_X(w http.ResponseWriter, r *http.Request) {
	x := r.URL.Query().Get("x")

	fmt.Fprintf(w, "Сработал маршрут только для X!\nПараметр: x = %s\n(параметр y не передан)\n", x)
}

func main() {
	r := mux.NewRouter()

	// 1. Строгий маршрут (Два параметра)
	// Метод .Queries() заставляет роутер проверять наличие ключей "x" и "y".
	// Важное правило: более строгие (длинные) маршруты нужно регистрировать ПЕРВЫМИ.
	r.HandleFunc("/Q", handlerQ_XY).
		Methods(http.MethodGet).
		Queries("x", "{x}", "y", "{y}")

	// 2. Менее строгий маршрут (Один параметр)
	// Если клиент передал только x, запрос не пройдет в первый обработчик
	// и провалится сюда.
	r.HandleFunc("/Q", handlerQ_X).
		Methods(http.MethodGet).
		Queries("x", "{x}")

	log.Println("Server is running on port 3000")
	log.Println("Тест 1 (оба параметра):  GET http://localhost:3000/Q?x=10&y=20")
	log.Println("Тест 2 (один параметр):  GET http://localhost:3000/Q?x=10")
	log.Println("Тест 3 (без параметров): GET http://localhost:3000/Q (выдаст 404 Not Found)")

	// Запуск сервера с нашим роутером
	log.Fatal(http.ListenAndServe(":3000", r))
}
