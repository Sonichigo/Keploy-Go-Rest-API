// model.go

package model

import (
	"context"
	"database/sql"

	// tom: errors is removed once functions are implemented
	// "errors"
	"fmt"
)

// tom: add backticks to json
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// tom: these are added after tdd tests
func (p *Post) GetPost(ctx context.Context, db *sql.DB) error {
	return db.QueryRowContext(ctx, "SELECT title, content FROM posts WHERE id=$1",
		p.ID).Scan(&p.Title, &p.Content)
}

func (p *Post) UpdatePost(ctx context.Context, db *sql.DB) error {
	_, err :=
		db.ExecContext(ctx, "UPDATE posts SET title=$1, content=$2 WHERE id=$3",
			p.Title, p.Content, p.ID)

	return err
}

func (p *Post) DeletePost(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "DELETE FROM posts WHERE id=$1", p.ID)

	return err
}

func (p *Post) CreatePost(ctx context.Context, db *sql.DB) error {
	err := db.QueryRowContext(ctx,
		"INSERT INTO posts(title, content) VALUES($1, $2) RETURNING id",
		p.Title, p.Content).Scan(&p.ID)

	if err != nil {
		return err
	}
	fmt.Printf("%v\n", p.ID)
	return nil
}

func GetAllPosts(ctx context.Context, db *sql.DB, start, count int) ([]Post, error) {
	rows, err := db.QueryContext(ctx,
		"SELECT id, title,  content FROM posts LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Posts := []Post{}

	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			return nil, err
		}
		Posts = append(Posts, p)
	}

	return Posts, nil
}
