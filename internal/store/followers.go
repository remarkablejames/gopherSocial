package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (s FollowersStore) Follow(ctx context.Context, followerID, followeeID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	// Check if the follower relationship already exists
	checkQuery := `SELECT 1 FROM followers WHERE user_id=$1 AND follower_id=$2`
	var exists int
	err := s.db.QueryRowContext(ctx, checkQuery, followeeID, followerID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If the relationship exists, return nil
	if exists == 1 {
		return nil
	}

	// Insert the new follower relationship
	query := `INSERT INTO followers (user_id, follower_id) VALUES($1, $2)`
	_, err = s.db.ExecContext(ctx, query, followeeID, followerID)
	if err != nil {
		return err
	}
	return nil
}

func (s FollowersStore) Unfollow(ctx context.Context, followerID, followeeID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	query := `DELETE FROM followers WHERE user_id=$1 AND follower_id=$2`

	_, err := s.db.ExecContext(ctx, query, followeeID, followerID)

	if err != nil {
		return err
	}
	return nil
}
