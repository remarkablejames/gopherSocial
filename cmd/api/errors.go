package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("INTERNAL SERVER ERROR: %s %s %s", r.RemoteAddr, r.Method, r.URL)
	err = WriteJSONError(w, http.StatusInternalServerError, "the server encountered an error and could not process the request")
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("BAD REQUEST ERROR: %s %s %s", r.RemoteAddr, r.Method, r.URL)
	err = WriteJSONError(w, http.StatusBadRequest, err.Error())
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("NOT FOUND ERROR: %s %s %s", r.RemoteAddr, r.Method, r.URL)
	err := WriteJSONError(w, http.StatusNotFound, "the requested resource could not be found")
	if err != nil {
		log.Fatal(err)
	}
}
