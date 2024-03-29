package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

// find the type
// match with the interface
// unmarshal and call on shape service

// The different operations:
// - Change shape (Update)
// - Change color
// - Change size

type ShapeRepo struct {
	database *pgxpool.Pool
}


type ShapeObject struct {
	Id		string 	`json:"id"`
	Shape 	string  `json:"shape"`
	Color 	string  `json:"color"`
	Size 	string  `json:"size"`
	Counter uint64	`json:"counter"`
}

func NewShapeRepository(database *pgxpool.Pool) *ShapeRepo {
	return &ShapeRepo{database: database}
}

func (repo *ShapeRepo) CreateShape(uid string, shape string, color string, size string) *uint64 {
	query := `
	    INSERT INTO shape (uid, shape_id, color_id, shape_size, counter)
		VALUES (
			$1,
			(SELECT id FROM shape_definition
				WHERE sdef.shape_name = $2),
			(SELECT id FROM color_definition
				WHERE color_name = $3),
			$4,
			0
		)
		RETURNING counter; 
	`
	// Problem: will users need to make sure creating new object will return the object itself
	var counter uint64

	err := repo.database.QueryRow(
		context.Background(),
		query,
		uid,
		shape,
		color,
		size,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
		os.Exit(1)
	}
	
	return &counter
}

func (repo *ShapeRepo) DeleteShape(uid string) *uint64 {
	query := `
	    DELETE FROM shape
		WHERE uid = $1
		RETURNING shape.counter + 1
	`
	var counter uint64

	err := repo.database.QueryRow(
		context.Background(),
		query,
		uid,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows are returned for query [%s\n]", query)
	}

	return &counter
}

func (repo *ShapeRepo) GetShape() *ShapeObject {
	query := `
		SELECT uid, sdef.shape_name, cdef.color_name, shape_size, counter 
		FROM shape INNER JOIN shape_definition AS sdef 
		ON shape.shape_id = sdef.id
		INNER JOIN color_definition AS cdef
		ON shape.color_id = cdef.id`
	
	var shape ShapeObject

	err := repo.database.QueryRow(
		context.Background(),
		query,
	).Scan(&shape.Id, &shape.Shape, &shape.Color, &shape.Size, &shape.Counter)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows are returned for query [%s\n]", query)
		log.Println(err);
	}

	return &shape
}	

func (repo *ShapeRepo) SetShape(uid string, shape string) *uint64 {
	query := `
	UPDATE shape 
	SET shape_id = (
		SELECT id FROM shape_definition 
		WHERE shape_name = $2
	),
	counter = counter + 1
	WHERE uid = $1
	RETURNING counter;`
		
    var counter uint64

	err := repo.database.QueryRow(
		context.Background(),
		query,
		uid,
		shape,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
		os.Exit(1)
	}

	return &counter
}

func (repo *ShapeRepo) SetColor(uid string, color string) *uint64 {
	query := `
	UPDATE shape 
	SET color_id = (
		SELECT id FROM color_definition 
		WHERE color_name = $2
	),
	counter = counter + 1
	WHERE uid = $1
	RETURNING counter`
	
    var counter uint64

	err := repo.database.QueryRow(
		context.Background(),
		query,
		uid,
		color,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
	}

	return &counter
}

func (repo *ShapeRepo) SetSize(uid string, size string) *uint64 {
	query := `
	UPDATE shape 
	SET shape_size = $2,
	counter = counter + 1
	WHERE uid = $1
	RETURNING counter`
	// TODO: separate counter into a different table
    var counter uint64

	err := repo.database.QueryRow(
		context.Background(),
		query,
		uid,
		size,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
	}

	return &counter
}