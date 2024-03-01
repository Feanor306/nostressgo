package client

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
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
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
			}
		case msg := <-h.Broadcast:
			fmt.Println(msg)
			// ONLY BROADCAST FOR CLIENTS THAT HAVE MATCHING SUBSCRIPTION ID
			// for client := range h.Clients {
			// 	select {
			// 	case client.Send <- message:
			// 	default:
			// 		close(client.Send)
			// 		delete(h.Clients, client)
			// 	}
			// }

			// for client, _ := range h.Clients {
			//     if err := client.Conn.WriteJSON(message); err != nil {
			//         fmt.Println(err)
			//         return
			//     }
			// }
		}

	}
}
