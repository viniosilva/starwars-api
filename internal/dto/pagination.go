package dto

type Pagination struct {
	Count    int    `json:"count" example:"10"`
	Total    int    `json:"total" example:"60"`
	Previous string `json:"previous" example:"http://localhost:8080/api/planets?size=10"`
	Next     string `json:"next" example:"http://localhost:8080/api/planets?page=3&size=10"`
}
