package tcp

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

const (
	DEFAULT_BUFFER_SIZE = 1024
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  DEFAULT_BUFFER_SIZE,
	WriteBufferSize: DEFAULT_BUFFER_SIZE,
}

type Server struct {
	port     string
	logInfo  *log.Logger
	logError *log.Logger
	upgrader *websocket.Upgrader
	hub      *Hub
}

func MakeServer(port string) *Server {
	return &Server{
		port:     port,
		logInfo:  log.New(os.Stdout, "INFO:Server:", log.Ldate|log.Ltime),
		logError: log.New(os.Stdout, "ERROR:Server:", log.Ldate|log.Ltime),
		upgrader: &Upgrader,
		hub:      MakeHub(),
	}
}

func (s *Server) Run() {
	go s.hub.Run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.serveWs(w, r)
	})
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		s.logError.Fatal("ListenAndServe: ", err)
	}
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logError.Println("serveWs: ", err)
		return
	}
	client := MakeClient(s.hub, conn)
	client.Hub.Register(client)
	go client.WritePump()
	go client.ReadPump()
}
