// client.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := "ws://localhost:3001/ws"
	fmt.Printf("Попытка подключения к %s...\n", serverURL)

	// Устанавливаем соединение с сервером
	c, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Ошибка подключения к серверу:", err)
	}
	// Закрываем соединение по завершении работы
	defer c.Close()

	fmt.Println("Подключение установлено!")

	// Текст сообщения для отправки
	message := "Привет, Gorilla WebSocket!"
	
	// Отправляем сообщение
	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
		return
	}
	fmt.Printf("Отправлено: %s\n", message)

	// Читаем ответ от сервера
	// Устанавливаем небольшой таймаут для чтения, чтобы клиент не завис навсегда, если сервер молчит
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, response, err := c.ReadMessage()
	if err != nil {
		log.Println("Ошибка при чтении ответа:", err)
		return
	}
	
	fmt.Printf("Получен ответ: %s\n", response)
}