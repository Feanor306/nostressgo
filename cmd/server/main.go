package main

import (
	"fmt"
	"net/http"

	"github.com/feanor306/nostressgo/src/config"
	"github.com/feanor306/nostressgo/src/database"
	"github.com/feanor306/nostressgo/src/handlers"
	"github.com/feanor306/nostressgo/src/server"
	"github.com/feanor306/nostressgo/src/service"
)

func serveWs(w http.ResponseWriter, r *http.Request, svc *service.Service) {
	fmt.Println("WebSocket Endpoint Hit")

	conn, err := svc.Upgrader.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := server.NewClient(conn, svc)

	client.Read()
}

func setupRoutes(svc *service.Service) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("NOSTRess GO started on port %d", svc.Conf.Port))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, svc)
	})
}

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	connString, err := conf.GetPostgresConnString()
	if err != nil {
		panic(err)
	}

	db, err := database.NewDatabase(connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.InitDatabase()
	if err != nil {
		panic(err)
	}

	svc := service.NewService(conf, db, handlers.NewUpgrader())

	fmt.Println("NOSTRess go!")
	setupRoutes(svc)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}