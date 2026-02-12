package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func InitDB(host string, port int)*sql.DB{
	connStr := "postgres://postgres:pass@" + host + fmt.Sprint(port) + "/postgres?sslmode=disable"
	db,err := sql.Open("postgres",connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

