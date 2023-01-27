package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type Server struct {
	handler *OperationManager
	config *Config
	hub *Hub
}

type ShapeOperations struct {
	OpType 		string		`json:"opType"`
	UuId		string		`json:"uuId"`
	ConflictId 	string  	`json:"conflictId"`
	IsDeleted 	bool		`json:"isDeleted"`
	Payload 	interface{}	`json:"payload"`
}

type AckOperations struct {
	Status 		string 		`json:"status"`
	OpType		string		`json:"opType"`
	UuId		string		`json:"uuId"`
	ConflictId 	string		`json:"conflictId"`
	Payload 	interface{} `json:"payload"`
}

type Payload struct {
	Id string 		  `json:"id"`
	NewShape string	`json:"newShape"`
	NewColor string `json:"newColor"`
	NewSize string  `json:"newSize"`
	NewCounter uint64 `json:"newCounter"`
	Shape 	string  `json:"shape"`
	Color 	string  `json:"color"`
	Size 	string  `json:"size"`
	Counter uint64	`json:"counter"`
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
	uuid := uuid.NewString()
	client := NewClient(s.hub, conn, s.handler, uuid) // possibly store sessionId here
	client.hub.register <- client
	s.handler.initInfo(conn, uuid) // not so neat, as it is not tightly coupled to client
	// problem: the uuid is sent to the user the first time, but not properly stored.
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	log.Println("Client Connected")
}

