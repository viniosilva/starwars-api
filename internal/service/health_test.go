package service_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_HealthService_Ping(t *testing.T) {
	var cases = map[string]struct {
		mocking     func(db sqlmock.Sqlmock)
		expectedErr error
	}{
		"should ping": {
			mocking: func(db sqlmock.Sqlmock) {
				db.ExpectPing()
			},
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

			healthService := service.IHealthService{DB: db}

			cs.mocking(mockDB)

			// when
			err = healthService.Ping(context.Background())

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
