package service

import (
	"context"

	"github.com/viniosilva/starwars-api/internal/model"
)

//go:generate mockgen -destination=../../mock/film_service_mock.go -package=mock . FilmService
type FilmService interface {
	CreateFilms(ctx context.Context, Films []model.Film) error
}

type IFilmService struct{}

func (impl *IFilmService) CreateFilms(ctx context.Context, Films []model.Film) error {
	// TODO
	return nil
}
