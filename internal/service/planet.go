package service

import (
	"context"

	"github.com/viniosilva/starwars-api/internal/model"
)

//go:generate mockgen -destination=../../mock/planet_service_mock.go -package=mock . PlanetService
type PlanetService interface {
	CreatePlanets(ctx context.Context, planets []model.PlanetWithFilms) error
	FindPlanets(ctx context.Context, page, size int) ([]model.PlanetWithFilms, int, error)
	FindPlanetByID(ctx context.Context, planetID int) (*model.PlanetWithFilms, error)
	DeletePlanet(ctx context.Context, planetID int) error
}

type IPlanetService struct{}

func (impl *IPlanetService) CreatePlanets(ctx context.Context, planets []model.PlanetWithFilms) error {
	// TODO
	return nil
}

func (impl *IPlanetService) FindPlanets(ctx context.Context, page, size int) ([]model.PlanetWithFilms, int, error) {
	// TODO
	return nil, 0, nil
}

func (impl *IPlanetService) FindPlanetByID(ctx context.Context, planetID int) (*model.PlanetWithFilms, error) {
	// TODO
	return nil, nil
}

func (impl *IPlanetService) DeletePlanet(ctx context.Context, planetID int) error {
	//TODO
	return nil
}
