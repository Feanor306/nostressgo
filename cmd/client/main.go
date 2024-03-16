package main

import (
	"bytes"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/feanor306/nostressgo/src/logger"
	"github.com/feanor306/nostressgo/test/client"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:3000", "http service address")

func main() {
	flag.Parse()
	log := logger.New()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("dial")
	}
	defer c.Close()

	cl := client.NewClient()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msg("read error")
			}

			if len(message) == 0 {
				continue
			}

			tc := cl.GetTestCase(true)
			tcb, err := tc.SerializeResponse()
			if err != nil {
				log.Error().Err(err).Msg("serialize test error")
			}

			if bytes.Equal(tcb, message) {
				log.Info().Str("payload", string(message)).Msg("PASS")
			} else {
				log.Error().Str("expected", string(tcb)).Str("received", string(message)).Msg("FAIL")
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tc := cl.GetTestCase(false)
			tcb, err := tc.SerializeRequest()
			if err != nil {
				log.Error().Err(err).Msg("serialize test error")
				return
			}

			err = c.WriteMessage(websocket.TextMessage, tcb)
			if err != nil {
				log.Error().Err(err).Msg("write err")
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error().Err(err).Msg("write close")
				return
			}
			return
		}
	}
}
