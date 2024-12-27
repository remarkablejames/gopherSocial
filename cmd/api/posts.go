package main

import (
	"gopherSocial/internal/store"
	"log"
	"net/http"
)

type CreatePostInput struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var inputPost CreatePostInput

	if err := ReadJSON(w, r, &inputPost); err != nil {
		err := WriteJSONError(w, http.StatusBadRequest, err.Error())
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}
	userId := 2
	post := &store.Post{
		Title:   inputPost.Title,
		Content: inputPost.Content,
		Tags:    inputPost.Tags,
		// TODO: get the user id from the request context after auth
		UserID: int64(userId),
	}
	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		err := WriteJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	if err := WriteJSON(w, http.StatusCreated, post); err != nil {
		log.Fatal(err)
		return
	}
}
