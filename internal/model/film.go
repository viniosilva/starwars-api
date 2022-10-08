package model

import "time"

type Film struct {
	ID          int       `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Title       string    `db:"title"`
	Director    string    `db:"director"`
	ReleaseDate time.Time `db:"release_date"`
}
