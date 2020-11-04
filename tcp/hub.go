package tcp

import (
	"log"
	"os"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logInfo    *log.Logger
	logError   *log.Logger
}

func MakeHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		logInfo:    log.New(os.Stdout, "INFO:Hub:", log.Ldate|log.Ltime),
		logError:   log.New(os.Stdout, "ERROR:Hub:", log.Ldate|log.Ltime),
	}
}

func (h *Hub) Register(c *Client) {
	h.logInfo.Println("Registering client:", c.conn.RemoteAddr().String())
	h.register <- c
}

func (h *Hub) Unregister(c *Client) {
	h.logInfo.Println("Unregistering client:", c.conn.RemoteAddr().String())
	h.unregister <- c
}

func (h *Hub) Broadcast(message []byte) {
	h.logInfo.Println("Broadcasting message:", string(message))
	h.broadcast <- message
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
