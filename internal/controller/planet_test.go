package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/controller"
	"github.com/viniosilva/starwars-api/internal/dto"
	"github.com/viniosilva/starwars-api/internal/exception"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/mock"
)

func Test_PlanetController_FindPlanetsAndTotal(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(planetService *mock.MockPlanetService)
		inputPage          int
		inputSize          int
		inputLoadFilms     bool
		inputName          string
		expectedStatusCode int
		expectedBody       dto.PlanetsResponse
		expectedErr        dto.ApiError
	}{
		"should return planets list": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetsAndTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.FindPlanetsAndTotalResult{Count: 1, Total: 1, Next: false, Data: []*model.Planet{{ID: 1}}}, nil)
			},
			inputLoadFilms:     true,
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.PlanetsResponse{
				Pagination: dto.Pagination{
					Count: 1,
					Total: 1,
				},
				Data: []dto.PlanetDto{{
					ID:        1,
					CreatedAt: "0001-01-01 00:00:00",
					UpdatedAt: "0001-01-01 00:00:00",
				}},
			},
		},
		"should return planets list when there are pages": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetsAndTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.FindPlanetsAndTotalResult{Count: 1, Total: 3, Next: true, Data: []*model.Planet{{ID: 1}}}, nil)
			},
			inputPage:          2,
			inputSize:          1,
			inputLoadFilms:     true,
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.PlanetsResponse{
				Pagination: dto.Pagination{
					Count:    1,
					Total:    3,
					Previous: "localhost?page=1&size=1",
					Next:     "localhost?page=3&size=1",
				},
				Data: []dto.PlanetDto{{
					ID:        1,
					CreatedAt: "0001-01-01 00:00:00",
					UpdatedAt: "0001-01-01 00:00:00",
				}},
			},
		},
		"should return planets list when name is tatooine": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetsAndTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.FindPlanetsAndTotalResult{Count: 1, Total: 1, Next: false, Data: []*model.Planet{{ID: 1, Name: "Tatooine"}}}, nil)
			},
			inputName:          "tatooine",
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.PlanetsResponse{
				Pagination: dto.Pagination{
					Count: 1,
					Total: 1,
				},
				Data: []dto.PlanetDto{{
					ID:        1,
					Name:      "Tatooine",
					CreatedAt: "0001-01-01 00:00:00",
					UpdatedAt: "0001-01-01 00:00:00",
				}},
			},
		},
		"should return empty planets list": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetsAndTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.FindPlanetsAndTotalResult{}, nil)
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
		"should throw internal server error": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetsAndTotal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.FindPlanetsAndTotalResult{}, fmt.Errorf("error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        dto.ApiError{Error: "internal server error"},
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

			q := fmt.Sprintf("?page=%d&size=%d", cs.inputPage, cs.inputSize)
			if cs.inputLoadFilms {
				q += "&loadFilms=true"
			}
			if cs.inputName != "" {
				q += fmt.Sprintf("&name=%s", cs.inputName)
			}

			ctx.Request = httptest.NewRequest("GET", "/api/planets"+q, nil)
			mockPlanetService := mock.NewMockPlanetService(ctrl)

			planetController := &controller.IPlanetController{Host: "localhost", PlanetService: mockPlanetService}
			planetController.Configure(r.Group("/api"))

			cs.mocking(mockPlanetService)

			// when
			planetController.FindPlanetsAndTotal(ctx)

			var body dto.PlanetsResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			var bodyErr dto.ApiError
			json.Unmarshal(res.Body.Bytes(), &bodyErr)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedBody, body)
			assert.Equal(t, cs.expectedErr, bodyErr)
		})
	}
}

func Test_PlanetController_FindPlanetByID(t *testing.T) {
	var cases = map[string]struct {
		mocking            func(planetService *mock.MockPlanetService)
		inputPlanetID      int
		inputLoadFilms     bool
		expectedStatusCode int
		expectedBody       dto.PlanetResponse
		expectedErr        dto.ApiError
	}{
		"should return planet": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&model.Planet{ID: 1}, nil)
			},
			inputPlanetID:      1,
			inputLoadFilms:     true,
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.PlanetResponse{
				Data: dto.PlanetDto{
					ID:        1,
					CreatedAt: "0001-01-01 00:00:00",
					UpdatedAt: "0001-01-01 00:00:00",
				},
			},
		},
		"should throw bad request when planetID is invalid": {
			mocking:            func(planetService *mock.MockPlanetService) {},
			inputPlanetID:      0,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr: dto.ApiError{
				Error: "invalid planet id",
			},
		},
		"should throw not found": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &exception.NotFoundException{Message: "planet 1 not found"})
			},
			inputPlanetID:      1,
			expectedStatusCode: http.StatusNotFound,
			expectedErr: dto.ApiError{
				Error: "planet 1 not found",
			},
		},
		"should throw internal server error when find planet by id": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().FindPlanetByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			inputPlanetID:      1,
			inputLoadFilms:     true,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        dto.ApiError{Error: "internal server error"},
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

			q := ""
			if cs.inputLoadFilms {
				q = "?loadFilms=true"
			}

			ctx.Params = append(ctx.Params, gin.Param{Key: "planetID", Value: fmt.Sprint(cs.inputPlanetID)})
			ctx.Request = httptest.NewRequest("GET", "/api/planets"+q, nil)

			mockPlanetService := mock.NewMockPlanetService(ctrl)

			planetController := &controller.IPlanetController{PlanetService: mockPlanetService}
			planetController.Configure(r.Group("/api"))

			cs.mocking(mockPlanetService)

			// when
			planetController.FindPlanetByID(ctx)

			var body dto.PlanetResponse
			json.Unmarshal(res.Body.Bytes(), &body)

			var bodyErr dto.ApiError
			json.Unmarshal(res.Body.Bytes(), &bodyErr)

			// then
			assert.Equal(t, cs.expectedStatusCode, res.Result().StatusCode)
			assert.Equal(t, cs.expectedBody, body)
			assert.Equal(t, cs.expectedErr, bodyErr)
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
		"should throw bad request when planetID is invalid": {
			mocking:            func(planetService *mock.MockPlanetService) {},
			inputPlanetID:      0,
			expectedStatusCode: http.StatusBadRequest,
		},
		"should throw internal server error when delete planet": {
			mocking: func(planetService *mock.MockPlanetService) {
				planetService.EXPECT().DeletePlanet(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			inputPlanetID:      1,
			expectedStatusCode: http.StatusInternalServerError,
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

func Test_PlanetController_ParsePlanetDto(t *testing.T) {
	climates, _ := json.Marshal([]string{"arid"})
	terrains, _ := json.Marshal([]string{"desert"})

	var cases = map[string]struct {
		inputPlanet *model.Planet
		expectedDto dto.PlanetDto
	}{
		"should return planetDTO": {
			inputPlanet: &model.Planet{
				ID:        1,
				CreatedAt: time.Date(2014, 12, 9, 13, 50, 49, 0, time.UTC),
				UpdatedAt: time.Date(2014, 12, 20, 20, 58, 18, 0, time.UTC),
				Name:      "Tatooine",
				Climates:  climates,
				Terrains:  terrains,
			},
			expectedDto: dto.PlanetDto{
				ID:        1,
				CreatedAt: "2014-12-09 13:50:49",
				UpdatedAt: "2014-12-20 20:58:18",
				Name:      "Tatooine",
				Climates:  []string{"arid"},
				Terrains:  []string{"desert"},
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			planetController := &controller.IPlanetController{}

			// when
			dto := planetController.ParsePlanetDto(cs.inputPlanet)

			// then
			assert.Equal(t, cs.expectedDto, dto)
		})
	}
}

func Test_PlanetController_ParseFilmDto(t *testing.T) {
	var cases = map[string]struct {
		inputFilm   *model.Film
		expectedDto dto.FilmDto
	}{
		"should return planetDTO": {
			inputFilm: &model.Film{
				ID:          1,
				CreatedAt:   time.Date(2014, 12, 10, 14, 23, 31, 0, time.UTC),
				UpdatedAt:   time.Date(2014, 12, 20, 19, 49, 45, 0, time.UTC),
				Title:       "A New Hope",
				Episode:     4,
				Director:    "George Lucas",
				ReleaseDate: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
			},
			expectedDto: dto.FilmDto{
				ID:          1,
				CreatedAt:   "2014-12-10 14:23:31",
				UpdatedAt:   "2014-12-20 19:49:45",
				Title:       "A New Hope",
				Episode:     4,
				Director:    "George Lucas",
				ReleaseDate: "1977-05-25",
			},
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			planetController := &controller.IPlanetController{}

			// when
			dto := planetController.ParseFilmDto(cs.inputFilm)

			// then
			assert.Equal(t, cs.expectedDto, dto)
		})
	}
}
