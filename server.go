package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type Server struct {
	handler *OperationManager
	config *Config
	hub *Hub
}

type ShapeOperations struct {
	OpType 		string		`json:"opType"`
	ConflictId 	string  	`json:"conflictId"`
	Payload 	interface{}	`json:"payload"`
}

type AckOperations struct {
	Status 		string 		`json:"status"`
	OpType		string		`json:"opType"`
	ConflictId 	string		`json:"conflictId"`
	Payload 	interface{} `json:"payload"`
}

type ShapePayload struct {
	Id string 		`json:"id"`
	NewShape string	`json:"newShape"`
	NewCounter int8 `json:"newCounter"`
}

type ColorPayload struct {
	Id string		`json:"id"`
	NewColor string `json:"newColor"`
	NewCounter int8 `json:"newCounter"`
}

type SizePayload struct {
	Id string		`json:"id"`
	NewSize string  `json:"newSize"`
	NewCounter int8 `json:"newCounter"`
}

var WsMessageType = 1

func NewServer(handler *OperationManager, config *Config, hub *Hub) *Server {
    return &Server{
		handler: handler,
		config: config,
		hub: hub,
	}
}

func (s *Server) Run() {
	fmt.Println("Hello WebSocket")
	s.Handler()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) Handler() {
	http.HandleFunc("/ws", s.wsEndpoint)
}

func (s *Server) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// There will be many instances of Client, so it shouldn't be in the DI container
	// Problem: should handler be instantiated once? not rly impt
	client := NewClient(s.hub, conn, s.handler)
	client.hub.register <- client
	s.handler.shape(conn)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	log.Println("Client Connected")
}

