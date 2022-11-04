package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/mitchellh/mapstructure"
	"github.com/gorilla/websocket"
)

type Server struct {
	config *Config
	shapeService *ShapeService
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
}

type ColorPayload struct {
	Id string		`json:"id"`
	NewColor string `json:"newColor"`
}

type SizePayload struct {
	Id string		`json:"id"`
	NewSize string  `json:"newSize"`
}

var WsMessageType = 1

func NewServer(config *Config, shapeService *ShapeService) *Server {
    return &Server{
		config: config,
		shapeService: shapeService,
	}
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

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	s.shape(ws)
// TODO replace reader with postgres persistence
	s.reader(ws)
}

func (s *Server) reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		var operation *ShapeOperations
 
		err = json.Unmarshal(p, &operation)

		if err != nil {
			log.Println(err)
			return
		}

		s.executeSetter(conn, operation)
		
	}
}

func (s *Server) shape(conn *websocket.Conn) {
	shape := s.shapeService.GetShape()

	packet := ShapeOperations{
		OpType: "load",
		ConflictId: "",
		Payload: shape,
	}
	bytes, _ := json.Marshal(packet)
	log.Println(packet)
	
	err := conn.WriteMessage(WsMessageType, bytes)

	if err != nil {
		log.Println(err)
	}
}

func (s *Server) setShape(conn *websocket.Conn, op *ShapeOperations) {
	var payload = ShapePayload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	returnShape := s.shapeService.SetShape(payload.NewShape)

	payload.NewShape = *returnShape

	result := AckOperations{
		Status: "success",
		OpType: op.OpType,
		ConflictId: op.ConflictId,
		Payload: payload,
	}
	bytes, _ := json.Marshal(result)

	err = conn.WriteMessage(WsMessageType, bytes)

	if err != nil {
		log.Println(err)
	}
}

func (s *Server) setColor(conn *websocket.Conn, op *ShapeOperations) {
	var payload = ColorPayload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	returnColor := s.shapeService.SetColor(payload.NewColor)

	payload.NewColor = *returnColor

	result := AckOperations{
		Status: "success",
		OpType: op.OpType,
		ConflictId: op.ConflictId,
		Payload: payload,
	}
	bytes, _ := json.Marshal(result)

	err = conn.WriteMessage(WsMessageType, bytes)

	if err != nil {
		log.Println(err)
	}
}

func (s *Server) setSize(conn *websocket.Conn, op *ShapeOperations) {
	var payload = SizePayload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	returnSize := s.shapeService.SetSize(payload.NewSize)
	
	payload.NewSize = *returnSize

	result := AckOperations{
		Status: "success",
		OpType: op.OpType,
		ConflictId: op.ConflictId,
		Payload: payload,
	}
	bytes, _ := json.Marshal(result)

	err = conn.WriteMessage(WsMessageType, bytes)

	if err != nil {
		log.Println(err)
	}
}

func (s *Server) Run() {
	fmt.Println("Hello WebSocket")
	s.Handler()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s *Server) executeSetter(conn *websocket.Conn, op *ShapeOperations) {
	switch op.OpType {
		case "SET_SHAPE":
			s.setShape(conn, op)
		case "SET_COLOR":
			s.setColor(conn, op)
		case "SET_SIZE":
			s.setSize(conn, op)
	}
}
