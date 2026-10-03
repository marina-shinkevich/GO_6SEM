package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:3000/ws"
	log.Printf("Connecting to %s", url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	defer conn.Close()

	log.Println("Connected")

	for i := 1; i <= 5; i++ {
		msg := fmt.Sprintf("message %d", i)
		log.Printf("Sending: %s", msg)

		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Fatalf("Write error: %v", err)
		}

		_, reply, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
		log.Printf("Received: %s", string(reply))

		time.Sleep(1 * time.Second)
	}

	// Закрываем соединение корректно
	err = conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"),
	)
	if err != nil {
		log.Printf("Close error: %v", err)
	}

	log.Println("Client done")
}
