package model

type SwapiFilm struct {
	Url         string `json:"url"`
	Created     string `json:"created"`
	Edited      string `json:"edited"`
	Title       string `json:"title"`
	EpisodeID   int    `json:"episode_id"`
	Director    string `json:"director"`
	ReleaseDate string `json:"release_date"`
}

type SwapiFilmsResponse struct {
	SwapiPaginateResponse
	Results []SwapiFilm `json:"results"`
}
