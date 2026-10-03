package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


func handlerA(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)

	name := vars["name"]

	fmt.Fprintf(w, "Привет, %s!\n", name)
}

// handlerC — Обработчик для POST-запроса с данными application/x-www-form-urlencoded
// Ожидает классическую отправку данных HTML-формы [1].
func handlerC(w http.ResponseWriter, r *http.Request) {
	// 1. Обязательный шаг: парсинг формы [1, 3].
	// Метод ParseForm() считывает тело запроса (body) и преобразует строку
	// вида "name=John&message=Hello" во внутреннюю структуру (словарь) языка Go.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "cannot parse form", http.StatusBadRequest)
		return
	}

	// 2. Получаем конкретные параметры из тела запроса по их строковым ключам [1, 3].
	nameParam := r.FormValue("name")
	messageParam := r.FormValue("message")

	// 3. Формируем ответ (в лекциях для этого часто используется map) [1].
	resp := map[string]string{
		"name":        nameParam,
		"message":     messageParam,
		"contentType": r.Header.Get("Content-Type"),
		"info":        fmt.Sprintf("%s %s", r.Method, r.RequestURI),
	}

	// 4. Отправляем полученные данные обратно клиенту в формате JSON,
	// предварительно установив правильный заголовок [1].
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Создаем экземпляр роутера gorilla/mux [1, 2].
	r := mux.NewRouter()

	// 1. Маршрутизация с ШАБЛОНОМ (Path Parameters).
	// {name} — это переменная в пути. Роутер перехватит любой текст после /A/ [2].
	r.HandleFunc("/A/{name}", handlerA).Methods(http.MethodGet)

	// 2. Маршрутизация для обработки BODY-параметров (отправка формы).
	// Привязываем обработку строго к POST-запросу [1].
	r.HandleFunc("/C", handlerC).Methods(http.MethodPost)

	log.Println("Server is running on 3000")
	// Запускаем сервер, передавая наш настроенный роутер 'r' [1, 2].
	if err := http.ListenAndServe("127.0.0.1:3000", r); err != nil {
		log.Fatal(err)
	}
}
