package models

import "time"

type Post struct {
	ID          int       `db:"id" json:"id"`
	Created     time.Time `db:"created" json:"created"`
	Modified    time.Time `db:"modified" json:"modified"`
	Slug        string    `db:"slug" json:"slug"`
	Title       string    `db:"title" json:"title"`
	Body        string    `db:"body" json:"body"`
	AuthorID    int       `db:"author_id" json:"author_id`
	Picture     string    `db:"picture" json:"picture"`
	Description string    `db:"description" json:"description"`
	Published   bool      `db:"published" json:"published"`
}
