package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type profileStruct struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	DOB time.Time `json:"dob"`
	Bio string `json:"bio"`
	Hobbies []string `json:"hobbies"`
	CreatedAt time.Time `json:"created_at"`
}

func handleError(err error){
	if err != nil {
		fmt.Println("Error encountered: ",err)
	}
}

func initPQ(host string, port int) *sql.DB{
	connStr := "postgres://postgres:pass@" + host + ":" + fmt.Sprint(port) + "/postgres?sslmode=disable"
	db,err := sql.Open("postgres",connStr)
	handleError(err)
	return db
}

func fetchProfileData(db *sql.DB, email string)(profileStruct,error){
	sql_query := `SELECT id,email,username,dob,bio,hobbies,created_at FROM profiles WHERE email = $1`
	row := db.QueryRow(sql_query,email)
	var profile profileStruct
	err := row.Scan(
		&profile.ID, 
		&profile.Email,
		&profile.Username, 
		&profile.DOB, 
		&profile.Bio,
		pq.Array(&profile.Hobbies),
		&profile.CreatedAt,
	)
	return profile,err
}
