package handler

import "log"

// Keeps the content and the destination of a websockets message
type message struct {
	// Websockets client endpoint id, eid follows usp specification.
	// This field is needed for us to know which agent or controller
	// the message is intended to be delivered to.
	eid  string
	data []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// Global hub instance
var hub *Hub

func InitHandlers() {
	hub = newHub()
	hub.run()
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// register new eid
			h.clients[client.eid] = client
			log.Printf("New client connected: %s", client.eid)
		case client := <-h.unregister:
			// verify if eid exists
			if _, ok := h.clients[client.eid]; ok {
				// delete eid form map of connections
				delete(h.clients, client.eid)
				// close client messages receiving channel
				close(client.send)
			}
			log.Println("Disconnected client", client.eid)
		case message := <-h.broadcast:
			// verify if eid exists
			if c, ok := h.clients[message.eid]; ok {
				select {
				// send message to receiver client
				case c.send <- message.data:
					log.Printf("Sent a message to %s", message.eid)
				default:
					// in case the message sending fails, close the client connection
					// because it means that the client is no longer active
					log.Printf("Failed to send a message to %s, disconnecting client...", message.eid)
					close(c.send)
					delete(h.clients, c.eid)
				}
			} else {
				log.Printf("Message receiver not found: %s", message.eid)
			}
		}
	}
}
