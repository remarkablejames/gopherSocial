package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	PostID  int64  `json:"post_id"`
	UserID  int64  `json:"user_id"`
}

type CommentStore struct {
	db *sql.DB
}

func (s CommentStore) Create(ctx context.Context, comment *Comment) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	query := `INSERT INTO comments (content,post_id,user_id)
	VALUES($1, $2, $3) RETURNING id`

	err := s.db.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.UserID).Scan(&comment.ID)

	if err != nil {
		return err
	}
	return nil
}

func (s CommentStore) GetPostByID(ctx context.Context, id int64) ([]*Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	query := `SELECT c.id, c.content, c.post_id, c.user_id FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id=$1 ORDER BY c.id DESC`
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*Comment{}

	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
