package main

import (
	"context"
	"log"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Username 	string
	Password 	string
	DbName		string
	Port		string
}

type Config struct {
	DatabaseConfig *PostgresConfig
	Port string
}



func NewConfig(dbConfig *PostgresConfig) *Config {
    return &Config{
		DatabaseConfig: dbConfig,
		Port: os.Getenv("PORT"),
	}
}

func NewPgConfig() *PostgresConfig {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("error loading .env file")
	}
	var username 	= os.Getenv("USERNAME")
    var password	= os.Getenv("PASSWORD")
	var port 		= os.Getenv("DATABASE_PORT")
	var dbName		= os.Getenv("DB_NAME")

	return &PostgresConfig{
		Username: username,
		Password: password,
		Port: port,
		DbName: dbName,
	}
}

func PostgresSetup(config *Config) *pgxpool.Pool {
	

	var url = fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s", 
		config.DatabaseConfig.Username, 
		config.DatabaseConfig.Password, 
		config.DatabaseConfig.Port, 
		config.DatabaseConfig.DbName,
	)

	pool, err := pgxpool.New(context.Background(), url)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return pool
}