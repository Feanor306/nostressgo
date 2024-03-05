package client

import (
	"github.com/feanor306/nostressgo/src/types"
	"github.com/feanor306/nostressgo/src/utils"
	"github.com/nbd-wtf/go-nostr"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan *nostr.Event

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *nostr.Event),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Start() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			delete(h.Clients, client)
		case event := <-h.Broadcast:
			for client := range h.Clients {
				// don't broadcast events to self
				if client.PubKey == event.PubKey {
					continue
				}
				for _, sub := range client.Subscriptions {
					for _, filter := range sub.Filters {
						if utils.EventMatchesFilter(event, &filter) {
							ee := nostr.EventEnvelope{
								Event: *event,
							}
							client.Respond(&types.EnvelopeWrapper{Envelope: &ee})
						}
					}
				}
			}
		}

	}
}
