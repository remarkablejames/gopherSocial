package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"gopherSocial/internal/store"
	"log"
	"net/http"
	"strconv"
)

type postKey string

const postCtxKey postKey = "post"

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
	post := app.getPostFromContext(r)

	comments, err := app.store.Comments.GetPostByID(r.Context(), post.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	app.jsonResponse(w, r, post, http.StatusOK)

}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()
	err = app.store.Posts.Delete(ctx, id)
	if err != nil {
		switch {

		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return

		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostInput struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := app.getPostFromContext(r)

	var updatePost UpdatePostInput
	if err := ReadJSON(w, r, &updatePost); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(updatePost); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if updatePost.Content != nil {
		post.Content = *updatePost.Content
	}
	if updatePost.Title != nil {
		post.Title = *updatePost.Title
	}

	err := app.store.Posts.Update(r.Context(), post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, r, post, http.StatusOK)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		if postID == "" {
			app.notFoundResponse(w, r)
			return
		}
		id, err := strconv.ParseInt(postID, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		post, err := app.store.Posts.GetByID(r.Context(), id)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), postCtxKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getPostFromContext(r *http.Request) *store.Post {
	post, ok := r.Context().Value(postCtxKey).(*store.Post)
	if !ok {
		return nil
	}
	return post
}
