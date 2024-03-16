package main

import (
	"fmt"
	"net/http"

	"github.com/feanor306/nostressgo/src/client"
	"github.com/feanor306/nostressgo/src/config"
	"github.com/feanor306/nostressgo/src/database"
	"github.com/feanor306/nostressgo/src/handlers"
	"github.com/feanor306/nostressgo/src/service"
)

func serveWs(w http.ResponseWriter, r *http.Request, svc *service.Service, hub *client.Hub) {
	conn, err := svc.Upgrader.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	cl := client.NewClient(conn, svc, hub)
	hub.Register <- cl
	go cl.Read()
}

func setupRoutes(svc *service.Service) {
	hub := client.NewHub()
	go hub.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, fmt.Sprintf("NOSTRess GO started on port %d", svc.Conf.Port))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, svc, hub)
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
