package models

import "time"

type Post struct {
	ID          int       `db:"id"`
	Created     time.Time `db:"created"`
	Modified    time.Time `db:"modified"`
	Slug        string    `db:"slug"`
	Title       string    `db:"title"`
	Body        string    `db:"body"`
	AuthorID    int       `db:"author_id"`
	Picture     string    `db:"picture"`
	Description string    `db:"description"`
	Published   bool      `db:"published"`
}
