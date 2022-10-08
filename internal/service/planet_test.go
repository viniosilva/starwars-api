package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_PlanetService_CreatePlanets(t *testing.T) {
	var cases = map[string]struct {
		inputPlanets  []model.PlanetWithFilms
		expectedError error
	}{
		"should create planets": {
			inputPlanets: []model.PlanetWithFilms{{
				Planet: model.Planet{
					ID:        1,
					CreatedAt: time.Date(2014, 12, 9, 13, 50, 49, 641000, time.UTC),
					UpdatedAt: time.Date(2014, 12, 20, 20, 58, 18, 411000, time.UTC),
					Name:      "Tatooine",
					Climates:  []string{"arid"},
					Terrains:  []string{"desert"},
				},
				Films: []model.Film{{ID: 1}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}},
			},
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			planetService := service.IPlanetService{}

			// when
			err := planetService.CreatePlanets(context.Background(), cs.inputPlanets)

			// then
			assert.Equal(t, cs.expectedError, err)
		})
	}
}

func Test_PlanetService_FindPlanets(t *testing.T) {
	var cases = map[string]struct {
		inputPage       int
		inputSize       int
		expectedPlanets []model.PlanetWithFilms
		expectedTotal   int
		expectedError   error
	}{
		"should return planets list": {},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			planetService := service.IPlanetService{}

			// when
			planets, total, err := planetService.FindPlanets(context.Background(), cs.inputPage, cs.inputSize)

			// then
			assert.Equal(t, cs.expectedPlanets, planets)
			assert.Equal(t, cs.expectedTotal, total)
			assert.Equal(t, cs.expectedError, err)
		})
	}
}

func Test_PlanetService_FindPlanetByID(t *testing.T) {
	var cases = map[string]struct {
		inputPlanetID  int
		expectedPlanet *model.PlanetWithFilms
		expectedError  error
	}{
		"should return planet": {},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			planetService := service.IPlanetService{}

			// when
			planet, err := planetService.FindPlanetByID(context.Background(), cs.inputPlanetID)

			// then
			assert.Equal(t, cs.expectedPlanet, planet)
			assert.Equal(t, cs.expectedError, err)
		})
	}
}

func Test_PlanetService_DeletePlanet(t *testing.T) {
	var cases = map[string]struct {
		inputPlanetID int
		expectedError error
	}{
		"should delete planet": {},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			planetService := service.IPlanetService{}

			// when
			err := planetService.DeletePlanet(context.Background(), cs.inputPlanetID)

			// then
			assert.Equal(t, cs.expectedError, err)
		})
	}
}
