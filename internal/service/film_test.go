package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_FilmService_CreateFilms(t *testing.T) {
	var cases = map[string]struct {
		mocking     func(db sqlmock.Sqlmock)
		inputFilms  []*model.Film
		expectedErr error
	}{
		"should create films": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectExec("INSERT IGNORE INTO").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			inputFilms: []*model.Film{{
				ID:          1,
				CreatedAt:   time.Date(2014, 12, 10, 14, 23, 31, 88000, time.UTC),
				UpdatedAt:   time.Date(2014, 12, 20, 19, 49, 45, 25600, time.UTC),
				Title:       "A New Hope, Episode 4",
				Director:    "George Lucas",
				ReleaseDate: time.Date(1977, 05, 25, 0, 0, 0, 0, time.UTC),
			},
			},
		},
		"should throw error when insert": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectExec("INSERT IGNORE INTO").WillReturnError(fmt.Errorf("error"))
			},
			inputFilms: []*model.Film{{
				ID:          1,
				CreatedAt:   time.Date(2014, 12, 10, 14, 23, 31, 88000, time.UTC),
				UpdatedAt:   time.Date(2014, 12, 20, 19, 49, 45, 25600, time.UTC),
				Title:       "A New Hope, Episode 4",
				Director:    "George Lucas",
				ReleaseDate: time.Date(1977, 05, 25, 0, 0, 0, 0, time.UTC),
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

			filmService := service.IFilmService{DB: db}

			cs.mocking(mockDB)

			// when
			err = filmService.CreateFilms(context.Background(), cs.inputFilms)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
