package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_FilmService_CreateFilms(t *testing.T) {
	var cases = map[string]struct {
		inputFilms    []model.Film
		expectedError error
	}{
		"should create films": {
			inputFilms: []model.Film{{
				ID:          1,
				CreatedAt:   time.Date(2014, 12, 10, 14, 23, 31, 88000, time.UTC),
				UpdatedAt:   time.Date(2014, 12, 20, 19, 49, 45, 25600, time.UTC),
				Title:       "A New Hope, Episode 4",
				Director:    "George Lucas",
				ReleaseDate: time.Date(1977, 05, 25, 0, 0, 0, 0, time.UTC),
			},
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			filmService := service.IFilmService{}

			// when
			err := filmService.CreateFilms(context.Background(), cs.inputFilms)

			// then
			assert.Equal(t, cs.expectedError, err)
		})
	}
}
