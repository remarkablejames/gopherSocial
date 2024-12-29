package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"gopherSocial/internal/store"
	"net/http"
	"strconv"
)

type userKey string

const userCtxKey userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	app.jsonResponse(w, r, user, http.StatusOK)

}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	followerUserId := user.ID
	var followeeUserId int64 = 2 // TODO: hardcoding for now, will get from the request context after auth

	ctx := r.Context()
	err := app.store.Followers.Follow(ctx, followerUserId, followeeUserId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	followerUserId := user.ID
	var followeeUserId int64 = 2 // TODO: hardcoding for now, will get from the request context after auth

	ctx := r.Context()
	err := app.store.Followers.Unfollow(ctx, followerUserId, followeeUserId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch err {
			case store.ErrRecordNotFound:
				app.notFoundResponse(w, r)
			default:
				app.internalServerError(w, r, err)

			}
		}

		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// utility function to get the user from the context
func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtxKey).(*store.User)
	return user
}
