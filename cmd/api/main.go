package main

import (
	"log"
)

func main() {

	cfg := config{
		addr: ":8081",
	}
	app := &application{
		config: cfg,
	}

	log.Fatal(app.run())
}
