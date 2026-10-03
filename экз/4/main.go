package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Настройка Upgrader для превращения HTTP-соединения в WebSocket
var upgrader = websocket.Upgrader{
	// Разрешаем запросы с любых источников (в продакшене лучше настроить проверку)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 1. Апгрейд HTTP-соединения до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при апгрейде соединения:", err)
		return
	}
	defer ws.Close()
	fmt.Println("Клиент успешно подключился!")

	// 2. Цикл для чтения и отправки сообщений (Эхо-сервер)
	for {
		// Чтение сообщения от клиента
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println("Клиент отключился или произошла ошибка:", err)
			break
		}

		// Вывод полученного сообщения в консоль сервера
		fmt.Printf("Получено: %s\n", string(p))

		// Отправка сообщения обратно клиенту (Эхо)
		err = ws.WriteMessage(messageType, p)
		if err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Сервер запущен на :3001/ws")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}