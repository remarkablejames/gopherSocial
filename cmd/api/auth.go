package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"gopherSocial/internal/store"
	"log"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("registerUserHandlers______=============================")
	var payload RegisterUserPayload
	if err := ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// Hash the password

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// generate a UUID string for the user

	plainToken := uuid.New().String()
	//log the token. This is for testing purposes only TODO: remove this
	log.Println("INVITATION TOKEN======>", plainToken)
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Create and store the user
	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrorDuplicateEmail:
			app.badRequestResponse(w, r, err)

		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	app.jsonResponse(w, r, nil, http.StatusOK)

}
