package model

type SwapiPlanet struct {
	Url     string   `json:"url"`
	Created string   `json:"created"`
	Edited  string   `json:"edited"`
	Name    string   `json:"name"`
	Climate string   `json:"climate"`
	Terrain string   `json:"terrain"`
	Films   []string `json:"films"`
}

type SwapiPlanetsResponse struct {
	SwapiPaginateResponse
	Results []SwapiPlanet `json:"results"`
}
