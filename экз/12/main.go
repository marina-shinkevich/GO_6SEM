package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux" 
)


type PersonInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}


type PersonOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// handlerPostPerson — обработчик, который принимает и Path-параметры, и JSON.
func handlerPostPerson(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем ШАБЛОН МАРШРУТА (переменную из URL).
	// Если URL был /person/123, то vars["id"] будет равно "123".
	vars := mux.Vars(r)
	id := vars["id"]

	var in PersonInput

	// 2. Читаем и декодируем JSON из тела запроса.
	// Используем потоковый декодер, передавая ему r.Body.
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Формируем объект для ответа.
	// Берем ID из URL-пути, а Name и Age — из раскодированного JSON-тела.
	out := PersonOutput{
		ID:   id,
		Name: in.Name,
		Age:  in.Age,
	}

	// 4. Устанавливаем правильный HTTP-заголовок для JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 5. Сериализуем и отправляем ответ клиенту.
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Создаем новый роутер gorilla/mux
	r := mux.NewRouter()

	// Регистрируем маршрут.
	// {id} — это шаблон маршрута.
	// .Methods(http.MethodPost) — жестко привязываем к POST-запросу.
	r.HandleFunc("/person/{id}", handlerPostPerson).Methods(http.MethodPost)

	log.Println("Server is running on port 3000")
	log.Println("Test with: POST http://localhost:3000/person/123")

	// Запускаем сервер, передавая наш кастомный роутер
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
