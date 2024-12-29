package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")
var queryTimeoutDuration = 5 * time.Second

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		Get(ctx context.Context, id int64) (*Post, error)
		GetByID(ctx context.Context, id int64) (*Post, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, post *Post) error
	}

	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
	}

	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetPostByID(ctx context.Context, id int64) ([]*Comment, error)
	}

	Followers interface {
		Follow(ctx context.Context, followerID, followeeID int64) error
		Unfollow(ctx context.Context, followerID, followeeID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &postsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowersStore{db},
	}
}
