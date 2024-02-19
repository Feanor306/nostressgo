package main

import (
	"fmt"
	"net/http"

	"github.com/feanor306/nostressgo/src/server"
)

func serveWs(hub *server.Hub, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	// create server earlier in main and pass here
	srv := server.NewServer()
	conn, err := srv.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &server.Client{
		Conn: conn,
		Hub:  hub,
	}

	hub.Register <- client
	client.Read()
}

func setupRoutes() {
	hub := server.NewHub()
	go hub.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "NOSTRess GO on port :8080")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
}

func main() {
	fmt.Println("NOSTRess go!")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
