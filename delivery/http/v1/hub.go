package v1

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan interface{}
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan interface{}),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Println("client unregistered")

		case message := <-h.broadcast:
			h.mu.Lock()
			activeClients := make([]*websocket.Conn, 0, len(h.clients))
			for c := range h.clients {
				activeClients = append(activeClients, c)
			}
			h.mu.Unlock()

			for _, client := range activeClients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("websocket write error: %v", err)
					client.Close()
					h.mu.Lock()
					delete(h.clients, client)
					h.mu.Unlock()
				}
			}
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	h.register <- conn
}

func (h *Hub) BroadcastUpdate(data interface{}) {
	h.broadcast <- data
}
