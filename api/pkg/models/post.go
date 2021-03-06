package models

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID          int            `db:"id" json:"id"`
	Created     time.Time      `db:"created" json:"created"`
	Modified    time.Time      `db:"modified" json:"modified"`
	Slug        string         `db:"slug" json:"slug"`
	Title       string         `db:"title" json:"title"`
	Body        string         `db:"body" json:"body"`
	AuthorID    int            `db:"author_id" json:"author_id"`
	AuthorEmail string         `db:"author_email" json:"author"`
	Picture     string         `db:"picture" json:"picture"`
	Description string         `db:"description" json:"description"`
	Published   bool           `db:"published" json:"published"`
	Tags        pq.StringArray `db:"tags" json:"tags"`
}
