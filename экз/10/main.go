package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv" 
	"github.com/gorilla/mux"
)


func handlerQ(w http.ResponseWriter, r *http.Request) {
	xStr := r.URL.Query().Get("x")
	yStr := r.URL.Query().Get("y")

	q := map[string]string{
		"x": xStr,
		"y": yStr,
	}

	// Вывод строковых значений клиенту
	fmt.Fprintf(w, "Строковые параметры: x = %s , y = %s\n", q["x"], q["y"])


	x, _ := strconv.ParseFloat(xStr, 64)
	y, _ := strconv.ParseFloat(yStr, 64)

	// Выводим отформатированные числа (2 знака после запятой)
	fmt.Fprintf(w, "Числовые параметры: x = %.2f, y = %.2f\n", x, y)
}

func main() {
	// Создаем новый роутер
	r := mux.NewRouter()

	// 2. СТРОГАЯ маршрутизация (Фишка gorilla/mux).
	// Метод .Queries() заставляет роутер проверять наличие параметров.
	// Этот маршрут сработает ТОЛЬКО если в URL есть ключи "x" и "y".
	r.HandleFunc("/Q", handlerQ).
		Methods(http.MethodGet).
		Queries("x", "{x}", "y", "{y}")

	log.Println("Server is running on port 3000")
	log.Println("Test with: GET http://localhost:3000/Q?x=5.5&y=10.2")

	// Запускаем сервер
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
