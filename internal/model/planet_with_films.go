package model

type PlanetWithFilms struct {
	Planet
	Films []Film `db:"films"`
}
