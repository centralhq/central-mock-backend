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
	Shape string  `json:"shape"`
	Color string  `json:"color"`
	Size string   `json:"size"`
}

func NewShapeRepository(database *pgxpool.Pool) *ShapeRepo {
	return &ShapeRepo{database: database}
}

func (repo *ShapeRepo) GetShape() *ShapeObject {
	query := `
		SELECT sdef.shape_name, cdef.color_name, shape_size 
		FROM shape INNER JOIN shape_definition AS sdef 
		ON shape.shape_id = sdef.id
		INNER JOIN color_definition AS cdef
		ON shape.color_id = cdef.id`
	
	var shape ShapeObject

	err := repo.database.QueryRow(
		context.Background(),
		query,
	).Scan(&shape.Shape, &shape.Color, &shape.Size)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows are returned for query [%s\n]", query)
	}

	return &shape
}	

func (repo *ShapeRepo) SetShape(shape string) *string {
	query := `
	UPDATE shape 
	SET shape_id = (
		SELECT id FROM shape_definition 
		WHERE shape_name = $1
	)
	RETURNING (
		SELECT shape_name
		FROM shape_definition AS def 
		WHERE def.id = shape_id
	)`
		
    var returnShape string

	err := repo.database.QueryRow(
		context.Background(),
		query,
		shape,
	).Scan(&returnShape)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
		os.Exit(1)
	}

	return &returnShape
}

func (repo *ShapeRepo) SetColor(color string) *string {
	query := `
	UPDATE shape 
	SET color_id = (
		SELECT id FROM color_definition 
		WHERE color_name = $1
	)
	RETURNING (
		SELECT color_name 
		FROM color_definition AS def 
		WHERE def.id = color_id
	)`
	
    var returnColor string

	err := repo.database.QueryRow(
		context.Background(),
		query,
		color,
	).Scan(&returnColor)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
	}

	return &returnColor
}

func (repo *ShapeRepo) SetSize(size string) *string {
	query := `
	UPDATE shape 
	SET shape_size = $1
	RETURNING shape_size`
	
    var returnSize string

	err := repo.database.QueryRow(
		context.Background(),
		query,
		size,
	).Scan(&returnSize)

	if err != nil {
		fmt.Fprintf(os.Stderr, "No rows returned for query [%s\n]", query)
	}

	return &returnSize
}