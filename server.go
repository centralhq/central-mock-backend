package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	Id string 		`json:"id"`
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

func (s *Server) jsonResponse(w http.ResponseWriter, res interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options:", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

func (s *Server) Handler() {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)

	r.Get("/init", s.initHandler)
}

func (s *Server) initHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	uuid := uuid.NewString()
	packet := s.handler.initInfo(uuid)

	s.jsonResponse(w, packet, http.StatusOK)
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
	// problem: the uuid is sent to the user the first time, but not properly stored.
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	log.Println("Client Connected")
}

