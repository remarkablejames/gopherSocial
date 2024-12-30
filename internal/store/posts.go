package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type postsStore struct {
	db *sql.DB
}

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	Version   int      `json:"version"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Comments  []*Comment
	User      User `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

func (s postsStore) Create(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	query := `INSERT INTO posts (content,title,user_id,tags)
	VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at `

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s postsStore) Get(ctx context.Context, id int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts WHERE id=$1`

	post := &Post{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s postsStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `SELECT id, content, title, user_id, tags, version, created_at, updated_at FROM posts WHERE id=$1`

	post := &Post{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.Version, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	return post, nil
}

func (s postsStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `DELETE FROM posts WHERE id=$1`

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (s postsStore) Update(ctx context.Context, post *Post) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	query := `UPDATE posts SET content=$1, title=$2, version = version + 1 WHERE id=$3 AND version=$4 RETURNING version`

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.ID, post.Version).Scan(&post.Version)

	if err != nil {
		switch {

		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err

		}
	}

	return nil
}

func (s postsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `SELECT p.id, p.content, p.title, p.user_id, p.tags, p.version, p.created_at, p.updated_at, u.username, COUNT(c.id) AS comments_count
FROM posts p
    LEFT JOIN comments c ON p.id = c.post_id
    LEFT JOIN users u ON p.user_id = u.id
JOIN followers f ON p.user_id = f.follower_id OR p.user_id = $1
WHERE f.user_id = $1 OR p.user_id = $1
GROUP BY p.id, u.username ORDER BY p.created_at ` + fq.Sort + `
                          LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []PostWithMetadata{}
	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.Version, &post.CreatedAt, &post.UpdatedAt, &post.User.Username, &post.CommentsCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err

	}

	return posts, nil
}
