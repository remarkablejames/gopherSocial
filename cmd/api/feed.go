package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// get the feed for the user
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(70))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// return the feed to the user
	app.jsonResponse(w, r, feed, http.StatusOK)

}
