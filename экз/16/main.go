package main

import (
	"fmt"
	"log"    
	"net/http" 
)


func trace(w http.ResponseWriter, r *http.Request) {

	log.Printf("Запрос: %s %s\n", r.Method, r.URL.Path)


	fmt.Fprintf(w, "Метод: %s, Путь: %s\n", r.Method, r.URL.Path)
}

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	trace(w, r)
}

func handleGetA(w http.ResponseWriter, r *http.Request) {
	trace(w, r)
}

func handleGetAB(w http.ResponseWriter, r *http.Request) {
	trace(w, r)
}

// router — наш собственный диспетчер (ручная маршрутизация)
func router(w http.ResponseWriter, r *http.Request) {
	// Первый уровень маршрутизации: фильтруем по HTTP-методу (GET, POST и т.д.)
	switch r.Method {

	case "GET":
		// Второй уровень: фильтруем по конкретному пути (URL)
		switch r.URL.Path {
		case "/":
			handleGetRoot(w, r)
		case "/A":
			handleGetA(w, r)
		case "/A/B":
			handleGetAB(w, r)
		default:
			http.NotFound(w, r)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// 1. Регистрируем наш единый роутер на корневой путь "/"
	// Теперь абсолютно все запросы к серверу будут попадать в функцию router
	http.HandleFunc("/", router)

	// 2. Выводим сообщение о старте в консоль с помощью пакета log
	log.Println("Сервер запущен: http://localhost:3000")

	// 3. Запускаем сервер.
	// Оборачиваем вызов в log.Fatal, чтобы программа немедленно завершилась с ошибкой,
	// если сервер не сможет запуститься (например, если порт 3000 уже занят).
	log.Fatal(http.ListenAndServe(":3000", nil))
}
