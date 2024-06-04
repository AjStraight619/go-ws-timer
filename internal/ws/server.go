package ws

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Server struct {
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan []byte
	Clients     map[*Client]bool
	Countdown   int
	Ticker      *time.Ticker
	clientsLock sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Clients:    make(map[*Client]bool),
		Countdown:  60,
	}
}

func (server *Server) Run() {
	for {
		select {
		case client := <-server.Register:
			server.clientsLock.Lock()
			server.Clients[client] = true
			server.clientsLock.Unlock()
			log.Println("Client registered")

		case client := <-server.Unregister:
			server.clientsLock.Lock()
			if _, ok := server.Clients[client]; ok {
				delete(server.Clients, client)
				close(client.Send)
				log.Println("Client unregistered")
			}
			server.clientsLock.Unlock()

		case message := <-server.Broadcast:
			server.clientsLock.Lock()
			for client := range server.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(server.Clients, client)
				}
			}
			server.clientsLock.Unlock()
		}
	}
}

func (server *Server) StartTimer() {
	server.Ticker = time.NewTicker(1 * time.Second)
	defer server.Ticker.Stop()

	for {
		select {
		case <-server.Ticker.C:
			if server.Countdown <= 0 {
				server.Broadcast <- []byte(`{"type": "timer", "data": 0}`)
				return
			}
			message := []byte(fmt.Sprintf(`{"type": "timer", "data": %d}`, server.Countdown))
			server.Broadcast <- message
			server.Countdown--
		}
	}
}
