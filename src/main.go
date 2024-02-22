package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/feanor306/nostressgo/src/config"
	"github.com/feanor306/nostressgo/src/database"
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

func setupRoutes(conf *config.Config) {
	hub := server.NewHub()
	go hub.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("NOSTRess GO started on port %d", conf.Port))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
}

func main() {
	ctx := context.Background()
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	connString, err := conf.GetPostgresConnString()
	if err != nil {
		panic(err)
	}

	db, err := database.NewDatabase(ctx, connString)
	if err != nil {
		panic(err)
	}

	str, err := db.InitDatabase(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(str)

	fmt.Println("NOSTRess go!")
	setupRoutes(conf)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
