package dto

type PlanetDto struct {
	ID        int       `json:"id" example:"1"`
	CreatedAt string    `json:"created_at,omitempty" example:"2014-12-09 13:50:49"`
	UpdatedAt string    `json:"updated_at,omitempty" example:"2014-12-20 20:58:18"`
	Films     []FilmDto `json:"films,omitempty"`
	Name      string    `json:"name,omitempty" example:"Tatooine"`
	Climates  []string  `json:"climates,omitempty" example:"arid"`
	Terrains  []string  `json:"terrains,omitempty" example:"desert"`
}

type PlanetResponse struct {
	Data PlanetDto `json:"data"`
}

type PlanetsResponse struct {
	Pagination
	Data []PlanetDto `json:"data"`
}
