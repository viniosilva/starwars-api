package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_HealthService_Ping(t *testing.T) {
	var cases = map[string]struct {
		expectedError error
	}{
		"should ping": {},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			healthService := service.IHealthService{}

			// when
			err := healthService.Ping(context.Background())

			// then
			assert.Equal(t, cs.expectedError, err)
		})
	}
}
