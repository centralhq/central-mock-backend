package main

import (
	"errors"
	"encoding/json"
	"log"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type OperationManager struct {
	shapeService *ShapeService
}

func NewOperationManager(shapeService *ShapeService) *OperationManager {
	return &OperationManager{
		shapeService: shapeService,
	}
 }
func (s *OperationManager) initInfo(conn *websocket.Conn) {
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

func (s *OperationManager) setShapeAndWrite(conn *websocket.Conn, op *ShapeOperations) {
	bytes, err := s.setShape(op)

	if err != nil {
		log.Println(err)
	}

	err = conn.WriteMessage(WsMessageType, bytes)

	if err != nil {
		log.Println(err)
	}
}

func (s *OperationManager) setShape(op *ShapeOperations) ([]byte, error) {
	var payload = ShapePayload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	newCounter := s.shapeService.SetShape(payload.NewShape)

	payload.NewCounter = *newCounter

	result := AckOperations{
		Status: "success",
		OpType: op.OpType,
		UuId: op.UuId,
		ConflictId: op.ConflictId,
		Payload: payload,
	}
	bytes, err := json.Marshal(result)

	return bytes, err
}

func (s *OperationManager) setColorAndWrite(conn *websocket.Conn, op *ShapeOperations) {
	bytes, err := s.setColor(op)

	if err != nil {
		log.Println(err)
	}

	err = conn.WriteMessage(WsMessageType, bytes)

	if err != nil {
		log.Println(err)
	}
}

func (s *OperationManager) setColor(op *ShapeOperations) ([]byte, error) {
	var payload = ColorPayload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	newCounter := s.shapeService.SetColor(payload.NewColor)

	payload.NewCounter = *newCounter

	result := AckOperations{
		Status: "success",
		OpType: op.OpType,
		UuId: op.UuId,
		ConflictId: op.ConflictId,
		Payload: payload,
	}
	bytes, err := json.Marshal(result)

	return bytes, err
}

func (s *OperationManager) setSizeAndWrite(conn *websocket.Conn, op *ShapeOperations) {
	bytes, err := s.setSize(op)

	if err != nil {
		log.Println(err)
	}

	err = conn.WriteMessage(WsMessageType, bytes) // Problem: You have to change the write message

	if err != nil {
		log.Println(err)
	}	
}

func (s *OperationManager) setSize(op *ShapeOperations) ([]byte, error) {
	var payload = SizePayload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		return nil, err
	}

	newCounter := s.shapeService.SetSize(payload.NewSize)
	
	payload.NewCounter = *newCounter

	result := AckOperations{
		Status: "success",
		OpType: op.OpType,
		UuId: op.UuId,
		ConflictId: op.ConflictId,
		Payload: payload,
	}
	bytes, err := json.Marshal(result)
	
	return bytes, err
}

func (s *OperationManager) executeSetterAndWrite(conn *websocket.Conn, op *ShapeOperations) {
	switch op.OpType {
		case "SET_SHAPE":
			s.setShapeAndWrite(conn, op)
		case "SET_COLOR":
			s.setColorAndWrite(conn, op)
		case "SET_SIZE":
			s.setSizeAndWrite(conn, op)
	}
}

func (s *OperationManager) executeSetter(op *ShapeOperations) ([]byte, error) {
	switch op.OpType {
		case "SET_SHAPE":
			return s.setShape(op)
		case "SET_COLOR":
			return s.setColor(op)
		case "SET_SIZE":
			return s.setSize(op)
		default:
			return nil, errors.New("OpType invalid")
	}
}
