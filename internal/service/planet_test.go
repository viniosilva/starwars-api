package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/dto"
	"github.com/viniosilva/starwars-api/internal/exception"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_PlanetService_CreatePlanets(t *testing.T) {
	climates, _ := json.Marshal([]string{"arid"})
	terrains, _ := json.Marshal([]string{"desert"})

	var cases = map[string]struct {
		mocking      func(db sqlmock.Sqlmock)
		inputPlanets []*model.Planet
		expectedErr  error
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
			expectedErr: fmt.Errorf("error"),
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
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PlanetService_CreateRelationshipFilmsToPlanets(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(db sqlmock.Sqlmock)
		inputRelationships map[int][]int
		expectedErr        error
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
			expectedErr:        fmt.Errorf("error"),
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
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PlanetService_FindPlanetsAndTotal(t *testing.T) {
	var cases = map[string]struct {
		mocking        func(db sqlmock.Sqlmock)
		inputPage      int
		inputSize      int
		inputLoadFilms bool
		expectedRes    dto.FindPlanetsAndTotalResult
		expectedErr    error
	}{
		"should return planets list": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{model.PlanetColumns.ID}).
					AddRow(1).AddRow(2))
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
				db.ExpectCommit()
			},
			inputPage: 1,
			inputSize: 1,
			expectedRes: dto.FindPlanetsAndTotalResult{
				Count: 1,
				Total: 2,
				Next:  true,
				Data:  []*model.Planet{{ID: 1}},
			},
		},
		"should return planets list when page is 2": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{model.PlanetColumns.ID}).
					AddRow(1).AddRow(2))
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
				db.ExpectCommit()
			},
			inputPage: 2,
			inputSize: 1,
			expectedRes: dto.FindPlanetsAndTotalResult{
				Count: 1,
				Total: 3,
				Next:  true,
				Data:  []*model.Planet{{ID: 1}},
			},
		},
		"should return planets empty list": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{model.PlanetColumns.ID}))
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				db.ExpectCommit()
			},
			inputPage:   1,
			inputSize:   1,
			expectedRes: dto.FindPlanetsAndTotalResult{},
		},
		"should throw error when begin tx": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectBegin().WillReturnError(fmt.Errorf("error"))
			},
			inputPage:   1,
			inputSize:   1,
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when planets all rollback": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("error"))
				db.ExpectRollback().WillReturnError(fmt.Errorf("error"))
			},
			inputPage:   1,
			inputSize:   1,
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when planets count rollback": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{model.PlanetColumns.ID}))
				db.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("error"))
				db.ExpectRollback().WillReturnError(fmt.Errorf("error"))
			},
			inputPage:   1,
			inputSize:   1,
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when commit": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectBegin()
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{model.PlanetColumns.ID}))
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				db.ExpectCommit().WillReturnError(fmt.Errorf("error"))
			},
			inputPage:   1,
			inputSize:   1,
			expectedErr: fmt.Errorf("error"),
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
			res, err := planetService.FindPlanetsAndTotal(context.Background(), cs.inputPage, cs.inputSize, cs.inputLoadFilms)

			// then
			assert.Equal(t, cs.expectedRes, res)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PlanetService_FindPlanetByID(t *testing.T) {
	var cases = map[string]struct {
		mocking        func(db sqlmock.Sqlmock)
		inputPlanetID  int
		inputLoadFilms bool
		expectedPlanet *model.Planet
		expectedErr    error
	}{
		"should return planet": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{model.PlanetColumns.ID}).AddRow(1))
			},
			inputPlanetID:  1,
			expectedPlanet: &model.Planet{ID: 1},
		},
		"should throw not found exception": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{model.PlanetColumns.ID}))
			},
			inputPlanetID: 1,
			expectedErr:   &exception.NotFoundException{Message: "planet 1 not found"},
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
			planet, err := planetService.FindPlanetByID(context.Background(), cs.inputPlanetID, cs.inputLoadFilms)

			// then
			assert.Equal(t, cs.expectedPlanet, planet)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_PlanetService_DeletePlanet(t *testing.T) {
	var cases = map[string]struct {
		mocking       func(db sqlmock.Sqlmock)
		inputPlanetID int
		expectedErr   error
	}{
		"should delete planet": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			inputPlanetID: 1,
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
			err = planetService.DeletePlanet(context.Background(), cs.inputPlanetID)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
