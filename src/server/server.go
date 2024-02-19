package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}
