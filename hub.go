// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Payload

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Payload),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) toBytes(payload *Payload) []byte {

	bytes, err := json.Marshal(payload)

	if err != nil {
		log.Println(err)
	}

	return bytes
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true //insert sessionId here
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if message.Uuid == client.uuid {
					bytes := h.toBytes(message)
					select {
					case client.send <- bytes:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}