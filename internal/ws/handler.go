package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var server = NewServer()

func init() {
	go server.Run()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	log.Println("User connected successfully")

	client := &Client{Conn: conn, Send: make(chan []byte)}
	server.Register <- client

	defer func() {
		server.Unregister <- client
		conn.Close()
	}()

	go client.writePump()
	client.readPump()
}
