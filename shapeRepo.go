package main

import (
	"context"
	"fmt"
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
	Shape 	string  `json:"shape"`
	Color 	string  `json:"color"`
	Size 	string  `json:"size"`
	Counter int8	`json:"counter"`
}

func NewShapeRepository(database *pgxpool.Pool) *ShapeRepo {
	return &ShapeRepo{database: database}
}

func (repo *ShapeRepo) GetShape() *ShapeObject {
	query := `
		SELECT sdef.shape_name, cdef.color_name, shape_size, counter 
		FROM shape INNER JOIN shape_definition AS sdef 
		ON shape.shape_id = sdef.id
		INNER JOIN color_definition AS cdef
		ON shape.color_id = cdef.id`
	
	var shape ShapeObject

	err := repo.database.QueryRow(
		context.Background(),
		query,
	).Scan(&shape.Shape, &shape.Color, &shape.Size, &shape.Counter)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows are returned for query [%s\n]", query)
	}

	return &shape
}	

func (repo *ShapeRepo) SetShape(shape string) *int8 {
	query := `
	UPDATE shape 
	SET shape_id = (
		SELECT id FROM shape_definition 
		WHERE shape_name = $1
	),
	SET counter = counter + 1,
	RETURNING counter;`
		
    var counter int8

	err := repo.database.QueryRow(
		context.Background(),
		query,
		shape,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
		os.Exit(1)
	}

	return &counter
}

func (repo *ShapeRepo) SetColor(color string) *int8 {
	query := `
	UPDATE shape 
	SET color_id = (
		SELECT id FROM color_definition 
		WHERE color_name = $1
	),
	counter = counter + 1,
	RETURNING (
		SELECT color_name 
		FROM color_definition AS def 
		WHERE def.id = color_id
	)`
	
    var counter int8

	err := repo.database.QueryRow(
		context.Background(),
		query,
		color,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
	}

	return &counter
}

func (repo *ShapeRepo) SetSize(size string) *int8 {
	query := `
	UPDATE shape 
	SET shape_size = $1,
	counter = counter + 1,
	RETURNING counter`
	// TODO: separate counter into a different table
    var counter int8

	err := repo.database.QueryRow(
		context.Background(),
		query,
		size,
	).Scan(&counter)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
	}

	return &counter
}