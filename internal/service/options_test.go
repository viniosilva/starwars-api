package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/service"
)

func Test_OptionService_GetOptionWhere(t *testing.T) {
	var cases = map[string]struct {
		inputOptions  []service.Option
		expectedQuery string
		expectedArg   interface{}
	}{
		"should return query": {
			inputOptions:  []service.Option{service.OptionWhere("name like ?", "test")},
			expectedQuery: "name like ?",
			expectedArg:   "test",
		},
		"should return empty query when option not exist": {
			inputOptions:  []service.Option{},
			expectedQuery: "",
			expectedArg:   "",
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// when
			query, arg := service.GetOptionWhere(cs.inputOptions)

			// then
			assert.Equal(t, cs.expectedQuery, query)
			assert.Equal(t, cs.expectedArg, arg)
		})
	}
}
