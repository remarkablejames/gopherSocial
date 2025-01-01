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
		GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error)
	}

	Users interface {
		Create(ctx context.Context, tx *sql.Tx, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(ctx context.Context, token string) error
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

func withTx(ctx context.Context, db *sql.DB, txFunc func(context.Context, *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = txFunc(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
