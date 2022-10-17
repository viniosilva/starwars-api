package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/controller"
	"github.com/viniosilva/starwars-api/internal/dto"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/mock"
)

func Test_PlanetController_FindPlanets(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(planetService *mock.MockPlanetService)
		expectedStatusCode int
		expectedBody       dto.PlanetsResponse
	}{
		"should return empty planets list": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanets(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]model.Planet{}, 0, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.PlanetsResponse{
				Pagination: dto.Pagination{
					Count:    0,
					Total:    0,
					Previous: "",
					Next:     "",
				},
				Data: []dto.PlanetDto{},
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			res := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(res)
			ctx.Request = httptest.NewRequest("GET", "/api/planets", nil)

			mockPlanetService := mock.NewMockPlanetService(ctrl)

			planetController := &controller.IPlanetController{PlanetService: mockPlanetService}
			planetController.Configure(r.Group("/api"))

			cs.mocking(mockPlanetService)

			// when
			planetController.FindPlanets(ctx)

			var body dto.PlanetsResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedBody, body)
		})
	}
}

func Test_PlanetController_FindPlanetByID(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(planetService *mock.MockPlanetService)
		inputPlanetID      int
		expectedStatusCode int
		expectedBody       dto.PlanetResponse
	}{
		"should return planet": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetByID(gomock.Any(), gomock.Any()).
					Return(&model.Planet{ID: 1}, nil)
			},
			inputPlanetID:      1,
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.PlanetResponse{
				Data: dto.PlanetDto{
					ID:        1,
					CreatedAt: "0001-01-01 00:00:00",
					UpdatedAt: "0001-01-01 00:00:00",
				},
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			res := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(res)
			ctx.Params = append(ctx.Params, gin.Param{Key: "planetID", Value: fmt.Sprint(cs.inputPlanetID)})
			ctx.Request = httptest.NewRequest("GET", "/api/planets", nil)

			mockPlanetService := mock.NewMockPlanetService(ctrl)

			planetController := &controller.IPlanetController{PlanetService: mockPlanetService}
			planetController.Configure(r.Group("/api"))

			cs.mocking(mockPlanetService)

			// when
			planetController.FindPlanetByID(ctx)

			var body dto.PlanetResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedBody, body)
		})
	}
}

func Test_PlanetController_DeletePlanet(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(planetService *mock.MockPlanetService)
		inputPlanetID      int
		expectedStatusCode int
	}{
		"should return planet": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().DeletePlanet(gomock.Any(), gomock.Any()).Return(nil)
			},
			inputPlanetID:      1,
			expectedStatusCode: http.StatusNoContent,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)
			res := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(res)
			ctx.Params = append(ctx.Params, gin.Param{Key: "planetID", Value: fmt.Sprint(cs.inputPlanetID)})
			ctx.Request = httptest.NewRequest("DELETE", "/api/planets", nil)

			mockPlanetService := mock.NewMockPlanetService(ctrl)

			planetController := &controller.IPlanetController{PlanetService: mockPlanetService}
			planetController.Configure(r.Group("/api"))

			cs.mocking(mockPlanetService)

			// when
			planetController.DeletePlanet(ctx)

			var body dto.PlanetResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
		})
	}
}
