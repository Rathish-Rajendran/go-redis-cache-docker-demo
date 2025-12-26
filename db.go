package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDbClient(){
	var err error
	ctx := context.Background()
	DSN := os.Getenv("DATABASE_URL")
	// Try connecting multiple times because Postgres might not be ready yet
	for attempt := 1; attempt <= 10; attempt++ {
		DB, err = pgxpool.New(ctx, DSN)
		if err != nil {
			log.Printf("Attempt %d: failed to create pool: %v", attempt, err)
		} else {
			// Check if the DB is actually ready
			err = DB.Ping(ctx)
			if err == nil {
				log.Println("Connected to PostgreSQL")
				break
			}
			log.Printf("Attempt %d: DB not ready yet: %v", attempt, err)
		}
		time.Sleep(3 * time.Second)
	}

	// Create table
	_, err = DB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			role TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Unable to create 'users' table:%v", err)
	}
	fmt.Println("Successfully setup database and created 'users' table")
}
