package server

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	ID  string
	Hub *Hub
	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send            chan []byte
	SubscriptionIds []string
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))
		var receivedEvent nostr.Event

		err = json.Unmarshal(p, &receivedEvent)
		if err != nil {
			log.Println(err)
		}

		// Should we broadcast on read?
		// c.Hub.Broadcast <- p
		// fmt.Printf("Message Received: %+v\n", receivedEvent.ID)
	}
}
