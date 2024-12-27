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
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var inputPost CreatePostInput

	if err := ReadJSON(w, r, &inputPost); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(inputPost); err != nil {
		app.badRequestResponse(w, r, err)
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
		app.internalServerError(w, r, err)
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
		app.internalServerError(w, r, err)
	}
	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	comments, err := app.store.Comments.GetPostByID(r.Context(), id)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := WriteJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
