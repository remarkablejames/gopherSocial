package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(`{"status": "ok"}`))
	if err := WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		log.Fatal(err)
	}
}
