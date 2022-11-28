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

func (s *OperationManager) toBytes(payload Payload, op *ShapeOperations) []byte {

	result := AckOperations{
		Status:     "success",
		OpType:     op.OpType,
		UuId:       op.UuId,
		ConflictId: op.ConflictId,
		Payload:    payload,
	}
	bytes, err := json.Marshal(result)

	if err != nil {
		log.Println(err)
	}

	return bytes
}

func (s *OperationManager) initInfo(conn *websocket.Conn, uuid string) {
	shape := s.shapeService.GetShape()

	packet := ShapeOperations{
		OpType: "load",
		UuId: uuid,
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

func (s *OperationManager) createShape(op *ShapeOperations) []byte {
	var payload = Payload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	newCounter := s.shapeService.CreateShape(
		payload.Id, 
		payload.Shape, 
		payload.Color,
		payload.Size,
	)

	payload.NewCounter = *newCounter

	return s.toBytes(payload, op)
}

func (s *OperationManager) deleteShape(op *ShapeOperations) []byte {
	var payload = Payload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	counter := s.shapeService.DeleteShape(payload.Id)

	/*
	Scenario 1: 
	- user 1 deletes from database,
	- user 2 modifies it after it's deletedd
	- user 2 receives an error, the rest doesn't
	Solving scenario 1: This means that user 2 will receive the delete, so it will be received.
	For now, remove user 2's edits.

	Scenario 2;
	- user 1 deletes from the database
	- user 2 modifies it before it's deleted
	- user 1 receives delete ack and so will user 2
	*/

	payload.NewCounter = *counter

	return s.toBytes(payload, op)
}

func (s *OperationManager) setShape(op *ShapeOperations) []byte {
	var payload = Payload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	newCounter := s.shapeService.SetShape(payload.Id, payload.NewShape)

	payload.NewCounter = *newCounter

	return s.toBytes(payload, op)
}

func (s *OperationManager) setColor(op *ShapeOperations) []byte {
	var payload = Payload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	newCounter := s.shapeService.SetColor(payload.Id, payload.NewColor)

	payload.NewCounter = *newCounter

	return s.toBytes(payload, op)
}

func (s *OperationManager) setSize(op *ShapeOperations) []byte {
	var payload = Payload{}
	err := mapstructure.Decode(op.Payload, &payload)

	if err != nil {
		log.Println(err)
	}

	newCounter := s.shapeService.SetSize(payload.Id, payload.NewSize)
	
	payload.NewCounter = *newCounter
	
	return s.toBytes(payload, op)
}

func (s *OperationManager) executeSetter(op *ShapeOperations) ([]byte, error) {
	switch op.OpType {
		case "DELETE_SHAPE":
			return s.deleteShape(op), nil

		case "CREATE_SHAPE":
			return s.createShape(op), nil
		case "SET_SHAPE":
			return s.setShape(op), nil
		case "SET_COLOR":
			return s.setColor(op), nil
		case "SET_SIZE":
			return s.setSize(op), nil
		default:
			return nil, errors.New("OpType invalid")
	}
}
