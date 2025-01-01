package main

import (
	"gopherSocial/internal/db"
	"gopherSocial/internal/store"
)

func main() {

	conn, err := db.New("postgres://gopheruser:123@localhost:5436/gophersocial?sslmode=disable", 10, 5, "5m")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	myStore := store.NewStorage(conn)
	db.Seed(myStore, conn)
}
