package client

import (
	"sync"

	"github.com/feanor306/nostressgo/src/types"
	"github.com/feanor306/nostressgo/src/utils"
	"github.com/nbd-wtf/go-nostr"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients.
	Clients      map[*Client]bool
	ClientsMutex sync.RWMutex

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
	go h.HandleRegister()
	go h.HandleUnregister()
	go h.HandleBroadcast()
}

func (h *Hub) HandleRegister() {
	for {
		client, ok := <-h.Register
		if !ok {
			break
		}
		h.ClientsMutex.Lock()
		h.Clients[client] = true
		h.ClientsMutex.Unlock()
	}
}

func (h *Hub) HandleUnregister() {
	for {
		client, ok := <-h.Unregister
		if !ok {
			break
		}
		h.ClientsMutex.Lock()
		delete(h.Clients, client)
		h.ClientsMutex.Unlock()
	}
}

func (h *Hub) HandleBroadcast() {
	for {
		event, ok := <-h.Broadcast
		if !ok {
			break
		}
		if len(h.Clients) == 0 {
			continue
		}
		h.ClientsMutex.Lock()
		for client := range h.Clients {
			// don't broadcast events to self
			if client.PubKey == event.PubKey || len(client.Subscriptions) == 0 {
				continue
			}
			for _, sub := range client.Subscriptions {
				if len(sub.Filters) == 0 {
					continue
				}
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
		h.ClientsMutex.Unlock()
	}
}
