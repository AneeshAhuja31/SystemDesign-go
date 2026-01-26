package node

import (
	"net/http"
	"database/sql"
	"fmt"
	"sha-go/ring"
)

type Node struct {
	ID string
	Slot int
	Hash uint64
	HttpServer *http.Server
	DB *sql.DB
}



func MakeNode(port int,httpaddr *http.Server,db *sql.DB,n int)*Node{
	id := "node-"+fmt.Sprint(port)
	hash := ring.Hash(id)
	slot := ring.GetSlot(hash,n)
	return &Node{
		ID: id,
		Slot: slot,
		Hash: hash,
		HttpServer: httpaddr,
		DB: db,
	}
}