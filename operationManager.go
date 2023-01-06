package main

import (
	"errors"
	"encoding/json"
	"reflect"
	"log"
	"github.com/gorilla/websocket"
)

// For vanilla, need to define a simpler struct and use that for the schema
type OperationManager struct {
	shapeService *ShapeService
}

func NewOperationManager(shapeService *ShapeService) *OperationManager {
	return &OperationManager{
		shapeService: shapeService,
	}
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

func (s *OperationManager) createShape(op *Payload) *Payload {

	newCounter := s.shapeService.CreateShape(
		op.Id, 
		op.Shape, 
		op.Color,
		op.Size,
	)

	op.NewCounter = *newCounter

	return op
}

func (s *OperationManager) deleteShape(op *Payload) *Payload {

	counter := s.shapeService.DeleteShape(op.Id)

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

	op.NewCounter = *counter

	return op
}

func (s *OperationManager) setShape(op *Payload) *Payload {

	newCounter := s.shapeService.SetShape(op.Id, op.NewShape)

	op.NewCounter = *newCounter

	return op
}

func (s *OperationManager) setColor(op *Payload) *Payload {

	newCounter := s.shapeService.SetColor(op.Id, op.NewColor)

	op.NewCounter = *newCounter

	return op
}

func (s *OperationManager) setSize(op *Payload) *Payload {

	newCounter := s.shapeService.SetSize(op.Id, op.NewSize)
	
	op.NewCounter = *newCounter
	
	return op
}

func (s *OperationManager) executeSetter(op *Payload) (*Payload, error) {
	t := reflect.TypeOf(*op)
	switch t.Field(2).Name {
		case "Shape":
			return s.createShape(op), nil
		case "NewShape":
			return s.setShape(op), nil
		case "NewColor":
			return s.setColor(op), nil
		case "NewSize":
			return s.setSize(op), nil
		default:
			return nil, errors.New("Payload type invalid")
	}
}