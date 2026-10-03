package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Настраиваем Upgrader, который "обновляет" обычное HTTP-соединение до WebSocket
var upgrader = websocket.Upgrader{
		// Разрешаем подключения с любых доменов (только для локальной разработки!)
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновляем первоначальный GET-запрос до WebSocket-соединения
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при обновлении соединения:", err)
		return
	}
	// Обязательно закрываем соединение при выходе из функции
	defer ws.Close()

	fmt.Println("Новый клиент успешно подключен!")

	// Бесконечный цикл для прослушивания новых сообщений
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Ошибка при чтении сообщения или клиент отключился:", err)
			break
		}
		fmt.Printf("Получено от клиента: %s\n", message)

		// Формируем ответ и отправляем его обратно клиенту
		response := []byte(fmt.Sprintf("Эхо от сервера: %s", message))
		err = ws.WriteMessage(messageType, response)
		if err != nil {
			log.Println("Ошибка при отправке ответа:", err)
			break
		}
	}
}

func main() {
	// Регистрируем обработчик для роута /ws
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("WebSocket сервер запущен на ws://localhost:3001/ws")
	// Запускаем HTTP-сервер
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}