package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type profile struct {
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

func fetchProfileData(db *sql.DB, email string)(profile,error){
	sql_query := `SELECT id,email,username,dob,bio,hobbies,created_at FROM profiles WHERE email = $1`
	row := db.QueryRow(sql_query,email)
	var user_profile profile
	err := row.Scan(
		&user_profile.ID, 
		&user_profile.Email,
		&user_profile.Username, 
		&user_profile.DOB, 
		&user_profile.Bio,
		pq.Array(&user_profile.Hobbies),
		&user_profile.CreatedAt,
	)
	return user_profile,err
}
