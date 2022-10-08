package exception_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/exception"
)

func Test_Exception_NotFoundException(t *testing.T) {
	var cases = map[string]struct {
		inputErrorMessage    string
		expectedErrorMessage string
	}{
		"should return error message": {
			inputErrorMessage:    "error",
			expectedErrorMessage: "error",
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// when
			error := exception.NotFoundException{Message: cs.inputErrorMessage}

			// then
			assert.Equal(t, cs.expectedErrorMessage, error.Error())
		})
	}
}
