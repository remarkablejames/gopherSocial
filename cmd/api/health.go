package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(`{"status": "ok"}`))
	app.jsonResponse(w, r, map[string]string{"status": "ok"}, http.StatusOK)
}
