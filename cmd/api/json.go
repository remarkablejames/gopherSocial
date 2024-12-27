package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// restrict the size of the request body to prevent abuse. 1MB is a reasonable size for most APIs.
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	// when sending error message, it is a good practice to wrap it in an envelope object and keep it consistent across the API.
	type envelope struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&envelope{Error: message})
}
