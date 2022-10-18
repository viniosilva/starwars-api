package dto

type FilmDto struct {
	ID          int    `json:"id" example:"1"`
	CreatedAt   string `json:"created_at,omitempty" example:"2014-12-09 13:50:49"`
	UpdatedAt   string `json:"updated_at,omitempty" example:"2014-12-20 20:58:18"`
	Title       string `json:"title,omitempty" example:"A New Hope"`
	Episode     int    `json:"episode,omitempty" example:"4"`
	Director    string `json:"director,omitempty" example:"George Lucas"`
	ReleaseDate string `json:"release_date,omitempty" example:"1977-05-25"`
}
