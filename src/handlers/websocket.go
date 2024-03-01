package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Upgrader struct {
	upgrader websocket.Upgrader
}

func NewUpgrader() *Upgrader {
	return &Upgrader{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := u.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}
