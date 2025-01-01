package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "error", err, "method", r.Method, "url", r.URL, "remote_addr", r.RemoteAddr)
	err = WriteJSONError(w, http.StatusInternalServerError, "the server encountered an error and could not process the request")
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "error", err, "method", r.Method, "url", r.URL, "remote_addr", r.RemoteAddr)
	err = WriteJSONError(w, http.StatusBadRequest, err.Error())
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Errorw("not found", "method", r.Method, "url", r.URL, "remote_addr", r.RemoteAddr)
	err := WriteJSONError(w, http.StatusNotFound, "the requested resource could not be found")
	if err != nil {
		log.Fatal(err)
	}
}
