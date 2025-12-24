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
			log.Println("New client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
				log.Println("Client unregistered and connection closed")
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			// Anti-constipation: مدیریت بهینه ارسال پیام
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket write error (broken pipe): %v", err)
					// Muscle Building: حذف کلاینت خراب بدون معطل کردن حلقه اصلی
					go func(c *websocket.Conn) {
						h.unregister <- c
					}(client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade failed: %v", err)
		return
	}
	h.register <- conn
}

func (h *Hub) BroadcastUpdate(data interface{}) {
	// استفاده از select برای جلوگیری از مسدود شدن در صورت پر بودن کانال
	select {
	case h.broadcast <- data:
	default:
		log.Println("Broadcast channel is full, skipping update")
	}
}
