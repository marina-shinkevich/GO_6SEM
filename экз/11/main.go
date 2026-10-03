package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv" 

	"github.com/gorilla/mux"
)


func handlerP1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]


	id, _ := strconv.Atoi(idStr)

	fmt.Fprintf(w, "id = %d\n", id)
}

// handlerP2 — обрабатывает маршрут сразу с двумя параметрами пути.
// Пример подходящего запроса: GET /P2/123/student
func handlerP2(w http.ResponseWriter, r *http.Request) {
	// Извлекаем все параметры маршрута.
	vars := mux.Vars(r)

	// Получаем значения по ключам, которые заданы в шаблоне роутера.
	idStr := vars["id"]
	sStr := vars["s"]

	// Конвертируем id в число (ошибку игнорируем для упрощения примера)
	id, _ := strconv.Atoi(idStr)

	// Формируем комбинированный ответ
	fmt.Fprintf(w, "id = %d, s = %s\n", id, sStr)
}

func main() {
	r := mux.NewRouter()

	
	r.HandleFunc("/P1/{id:[0-9]+}", handlerP1).Methods(http.MethodGet)

	// 2. Маршрут с несколькими параметрами.
	// Здесь мы комбинируем строгую проверку для {id} и любой текст для переменной {s}.
	r.HandleFunc("/P2/{id:[0-9]+}/{s}", handlerP2).Methods(http.MethodGet)
	log.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
