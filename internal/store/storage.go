package store

import (
	"context"
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		Get(ctx context.Context, id int64) (*Post, error)
		GetByID(ctx context.Context, id int64) (*Post, error)
	}

	Users interface {
		Create(ctx context.Context, user *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &postsStore{db},
		Users: &UsersStore{db},
	}
}
