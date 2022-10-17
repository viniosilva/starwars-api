package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_PlanetService_CreatePlanets(t *testing.T) {
	climates, _ := json.Marshal([]string{"arid"})
	terrains, _ := json.Marshal([]string{"desert"})

	var cases = map[string]struct {
		mocking       func(db sqlmock.Sqlmock)
		inputPlanets  []*model.Planet
		expectedError error
	}{
		"should create planets": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectExec("INSERT IGNORE INTO").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			inputPlanets: []*model.Planet{{
				ID:        1,
				CreatedAt: time.Date(2014, 12, 9, 13, 50, 49, 641000, time.UTC),
				UpdatedAt: time.Date(2014, 12, 20, 20, 58, 18, 411000, time.UTC),
				Name:      "Tatooine",
				Climates:  climates,
				Terrains:  terrains,
			}},
		},
		"should throw error when insert": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectExec("INSERT IGNORE INTO").WillReturnError(fmt.Errorf("error"))
			},
			inputPlanets: []*model.Planet{{
				ID:        1,
				CreatedAt: time.Date(2014, 12, 9, 13, 50, 49, 641000, time.UTC),
				UpdatedAt: time.Date(2014, 12, 20, 20, 58, 18, 411000, time.UTC),
				Name:      "Tatooine",
				Climates:  climates,
				Terrains:  terrains,
			}},
			expectedError: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			planetService := service.IPlanetService{DB: db}

			cs.mocking(mockDB)

			// when
			err = planetService.CreatePlanets(context.Background(), cs.inputPlanets)

			// then
			assert.Equal(t, cs.expectedError, err)
		})
	}
}

func Test_PlanetService_CreateRelationshipFilmsToPlanets(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(db sqlmock.Sqlmock)
		inputRelationships map[int][]int
		expectedError      error
	}{
		"should create relationships": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectExec("INSERT IGNORE INTO").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			inputRelationships: map[int][]int{1: {1}},
		},
		"should throw error when insert": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectExec("INSERT IGNORE INTO").WillReturnError(fmt.Errorf("error"))
			},
			inputRelationships: map[int][]int{1: {1}},
			expectedError:      fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			db, mockDB, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			planetService := service.IPlanetService{DB: db}

			cs.mocking(mockDB)

			// when
			err = planetService.CreateRelationshipFilmsToPlanets(context.Background(), cs.inputRelationships)

			// then
			assert.Equal(t, cs.expectedError, err)
		})
	}
}

func Test_PlanetService_FindPlanets(t *testing.T) {
	var cases = map[string]struct {
		inputPage       int
		inputSize       int
		expectedPlanets []model.Planet
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
		expectedPlanet *model.Planet
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
