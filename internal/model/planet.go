package model

import "time"

type Planet struct {
	ID        int        `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Name      string     `db:"name"`
	Climates  []string   `db:"climates"`
	Terrains  []string   `db:"terrains"`
}
