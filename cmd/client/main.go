package main

import (
	"bytes"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:3000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	client := newClient()

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatalf("read error: %s", err)
				return
			}
			if len(message) == 0 {
				continue
			}

			tc := client.GetTestCase(true)
			tcb, err := tc.SerializeResponse()
			if err != nil {
				log.Fatalf("serialize test error: %s", err)
				return
			}

			if bytes.Equal(tcb, message) {
				log.Printf("PASS received: %s", message)
			} else {
				log.Println("FAIL")
				log.Printf("expected: %s", tcb)
				log.Printf("received: %s", message)
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			tc := client.GetTestCase(false)
			tcb, err := tc.SerializeRequest()
			if err != nil {
				log.Fatalf("serialize test error: %s", err)
				return
			}

			err = c.WriteMessage(websocket.TextMessage, tcb)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
