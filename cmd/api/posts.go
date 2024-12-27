package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gopherSocial/internal/store"
	"log"
	"net/http"
	"strconv"
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

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request: %v\n", r)
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		err := WriteJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		err := WriteJSONError(w, http.StatusNotFound, err.Error())
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	if err := WriteJSON(w, http.StatusOK, post); err != nil {
		log.Fatal(err)
		return
	}
	// get the post id from the request path
	// get the post from the database
	// write the post as JSON to the response
}
