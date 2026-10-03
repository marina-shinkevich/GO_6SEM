package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected: %s", r.RemoteAddr)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Client disconnected normally: %s", r.RemoteAddr)
			} else {
				log.Printf("Client disconnected: %s (%v)", r.RemoteAddr, err)
			}
			return
		}

		log.Printf("Received from %s: %s", r.RemoteAddr, string(msg))

		reply := "from server: " + string(msg)
		if err := conn.WriteMessage(msgType, []byte(reply)); err != nil {
			log.Printf("Write error: %v", err)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", wsHandler)

	log.Println("WebSocket server starting on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
