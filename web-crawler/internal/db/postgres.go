package db

import (
	"os"
	"fmt"
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)


func InitDB()*sql.DB{
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB := os.Getenv("POSTGRES_DB")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		POSTGRES_USER,POSTGRES_PASSWORD,POSTGRES_HOST,POSTGRES_PORT,POSTGRES_DB)
	db,err := sql.Open("postgres",connStr)
	if err = db.Ping(); err != nil{
		log.Fatal("Postgres connection failed: ",err)
	}
	if err != nil{
		log.Fatal("Postgres connection failed with error: ",err)
	}

	initTableQuery := `
	CREATE TABLE IF NOT EXISTS frontier(
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL UNIQUE,
		domain TEXT NOT NULL,
		depth INT NOT NULL,
		status TEXT DEFAULT 'pending',
		priority INT DEFAULT 5,
		status_code INT,
		error_msg TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_frontier_dequeue ON frontier(status, priority, created_at);
	CREATE INDEX IF NOT EXISTS idx_frontier_domain ON frontier(domain);
	`
	_,e:= db.Exec(initTableQuery)
	if e != nil{
		log.Fatal("Error initializing database with error: ",e)
	}
	

	return db
}