package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux" 
)


type RequestD struct {
	Op string `json:"op"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

func handlerD(w http.ResponseWriter, r *http.Request) {
	var data RequestD

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Бизнес-логика: формируем ответ.
	// Используем map[string]interface{} — это универсальный словарь,
	// где ключи строковые, а значения могут быть ЛЮБОГО типа.
	resp := map[string]interface{}{
		"received":    data, // Вкладываем полученный объект целиком
		"op":          data.Op,
		"x":           data.X,
		"y":           data.Y,
		"contentType": r.Header.Get("Content-Type"),
		"info":        fmt.Sprintf("%s %s", r.Method, r.RequestURI),
	}

	// 3. Сообщаем клиенту, что отдаем данные в формате JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 4. Сериализуем наш словарь (map) в JSON и отправляем в поток ответа
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// 1. Создаем новый экземпляр роутера gorilla/mux
	r := mux.NewRouter()

	// 2. Регистрируем маршрут.
	// Главная фишка mux: мы можем сразу привязать обработчик строго к POST-методу
	// с помощью метода-цепочки .Methods(http.MethodPost)
	r.HandleFunc("/D", handlerD).Methods(http.MethodPost)

	log.Println("Server is running on", 3000)

	// 3. Запускаем сервер.
	// Важное отличие: вторым параметром передаем наш кастомный роутер 'r', а не 'nil'
	if err := http.ListenAndServe("127.0.0.1:3000", r); err != nil {
		log.Fatal(err)
	}
}
