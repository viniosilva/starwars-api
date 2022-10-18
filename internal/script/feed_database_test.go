package script_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/script"
	"github.com/viniosilva/starwars-api/mock"
)

func Test_FeedDatabaseScript_Execute(t *testing.T) {
	climates, _ := json.Marshal([]string{"arid"})
	terrains, _ := json.Marshal([]string{"desert"})

	swapiFilm := model.SwapiFilm{
		Url:         "https://swapi.dev/api/films/1/",
		Created:     "2014-12-10T14:23:31.880000Z",
		Edited:      "2014-12-20T19:49:45.256000Z",
		Title:       "A New Hope",
		EpisodeID:   4,
		Director:    "George Lucas",
		ReleaseDate: "1977-05-25",
	}
	swapiPlanet := model.SwapiPlanet{
		Url:     "https://swapi.dev/api/planets/1/",
		Created: "2014-12-09T13:50:49.641000Z",
		Edited:  "2014-12-20T20:58:18.411000Z",
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
		Films: []string{
			"https://swapi.dev/api/films/1/",
			"https://swapi.dev/api/films/3/",
			"https://swapi.dev/api/films/4/",
			"https://swapi.dev/api/films/5/",
			"https://swapi.dev/api/films/6/",
		},
	}
	film := model.Film{
		ID:          1,
		CreatedAt:   time.Date(2014, 12, 10, 14, 23, 31, 0, time.UTC),
		UpdatedAt:   time.Date(2014, 12, 20, 19, 49, 45, 0, time.UTC),
		Title:       "A New Hope",
		Episode:     4,
		Director:    "George Lucas",
		ReleaseDate: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
	}
	planet := model.Planet{
		ID:        1,
		CreatedAt: time.Date(2014, 12, 9, 13, 50, 49, 0, time.UTC),
		UpdatedAt: time.Date(2014, 12, 20, 20, 58, 18, 0, time.UTC),
		Name:      "Tatooine",
		Climates:  climates,
		Terrains:  terrains,
	}
	_, errInvalidUrl := strconv.Atoi("")

	var cases = map[string]struct {
		mocking     func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService)
		expectedErr error
	}{
		"should be successful": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(&model.SwapiPlanetsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiPlanet{swapiPlanet},
					}, nil)
				filmService.EXPECT().CreateFilms(gomock.Any(), []*model.Film{&film}).Return(nil)
				planetService.EXPECT().CreatePlanets(gomock.Any(), []*model.Planet{&planet}).Return(nil)
				planetService.EXPECT().CreateRelationshipFilmsToPlanets(gomock.Any(), map[int][]int{1: {1, 3, 4, 5, 6}}).Return(nil)
			},
		},
		"should throw error when get swapi films": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when get swapi planets": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when parse swapi film to model": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{{Url: ""}},
					}, nil)
			},
			expectedErr: errInvalidUrl,
		},
		"should throw error when parse swapi planet to model": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(&model.SwapiPlanetsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiPlanet{{Url: ""}},
					}, nil)
			},
			expectedErr: errInvalidUrl,
		},
		"should throw error when get id from url": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(&model.SwapiPlanetsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiPlanet{{Url: ""}},
					}, nil)
			},
			expectedErr: errInvalidUrl,
		},
		"should throw error when create films": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(&model.SwapiPlanetsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiPlanet{swapiPlanet},
					}, nil)
				filmService.EXPECT().CreateFilms(gomock.Any(), []*model.Film{&film}).Return(fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when create planets": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(&model.SwapiPlanetsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiPlanet{swapiPlanet},
					}, nil)
				filmService.EXPECT().CreateFilms(gomock.Any(), []*model.Film{&film}).Return(nil)
				planetService.EXPECT().CreatePlanets(gomock.Any(), []*model.Planet{&planet}).Return(fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when create relationship films to planets": {
			mocking: func(swapiRequest *mock.MockSwapiRequest, filmService *mock.MockFilmService, planetService *mock.MockPlanetService) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(&model.SwapiPlanetsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiPlanet{swapiPlanet},
					}, nil)
				filmService.EXPECT().CreateFilms(gomock.Any(), []*model.Film{&film}).Return(nil)
				planetService.EXPECT().CreatePlanets(gomock.Any(), []*model.Planet{&planet}).Return(nil)
				planetService.EXPECT().CreateRelationshipFilmsToPlanets(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

			},
			expectedErr: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSwapiRequest := mock.NewMockSwapiRequest(ctrl)
			mockFilmService := mock.NewMockFilmService(ctrl)
			mockPlanetService := mock.NewMockPlanetService(ctrl)
			feedDatabaseScript := &script.IFeedDatabaseScript{
				Swapi:         mockSwapiRequest,
				FilmService:   mockFilmService,
				PlanetService: mockPlanetService,
			}

			cs.mocking(mockSwapiRequest, mockFilmService, mockPlanetService)

			// when
			err := feedDatabaseScript.Execute()

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}

}

func Test_FeedDatabaseScript_GetSwapiFilms(t *testing.T) {
	swapiFilm := model.SwapiFilm{
		Url:         "https://swapi.dev/api/films/1/",
		Created:     "2014-12-10T14:23:31.880000Z",
		Edited:      "2014-12-20T19:49:45.256000Z",
		Title:       "A New Hope",
		EpisodeID:   4,
		Director:    "George Lucas",
		ReleaseDate: "1977-05-25",
	}

	var cases = map[string]struct {
		mocking            func(swapiRequest *mock.MockSwapiRequest)
		expectedSwapiFilms []model.SwapiFilm
		expectedErr        error
	}{
		"should return swapi films": {
			mocking: func(swapiRequest *mock.MockSwapiRequest) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(&model.SwapiFilmsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiFilm{swapiFilm},
					}, nil)
			},
			expectedSwapiFilms: []model.SwapiFilm{swapiFilm},
		},
		"should throw error when get films": {
			mocking: func(swapiRequest *mock.MockSwapiRequest) {
				swapiRequest.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSwapiRequest := mock.NewMockSwapiRequest(ctrl)
			feedDatabaseScript := &script.IFeedDatabaseScript{Swapi: mockSwapiRequest}

			cs.mocking(mockSwapiRequest)

			// when
			swapiPlanets, err := feedDatabaseScript.GetSwapiFilms(context.Background())

			// then
			assert.Equal(t, cs.expectedSwapiFilms, swapiPlanets)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FeedDatabaseScript_GetSwapiPlanets(t *testing.T) {
	swapiPlanet := model.SwapiPlanet{
		Url:     "https://swapi.dev/api/planets/1/",
		Created: "2014-12-09T13:50:49.641000Z",
		Edited:  "2014-12-20T20:58:18.411000Z",
		Name:    "Tatooine",
		Climate: "arid",
		Terrain: "desert",
		Films: []string{
			"https://swapi.dev/api/films/1/",
			"https://swapi.dev/api/films/3/",
			"https://swapi.dev/api/films/4/",
			"https://swapi.dev/api/films/5/",
			"https://swapi.dev/api/films/6/",
		},
	}

	var cases = map[string]struct {
		mocking              func(swapiRequest *mock.MockSwapiRequest)
		expectedSwapiPlanets []model.SwapiPlanet
		expectedErr          error
	}{
		"should return swapi planets": {
			mocking: func(swapiRequest *mock.MockSwapiRequest) {
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(&model.SwapiPlanetsResponse{
						SwapiPaginateResponse: model.SwapiPaginateResponse{Count: 1},
						Results:               []model.SwapiPlanet{swapiPlanet},
					}, nil)
			},
			expectedSwapiPlanets: []model.SwapiPlanet{swapiPlanet},
		},
		"should throw error when get planets": {
			mocking: func(swapiRequest *mock.MockSwapiRequest) {
				swapiRequest.EXPECT().GetPlanets(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSwapiRequest := mock.NewMockSwapiRequest(ctrl)
			feedDatabaseScript := &script.IFeedDatabaseScript{Swapi: mockSwapiRequest}

			cs.mocking(mockSwapiRequest)

			// when
			swapiPlanets, err := feedDatabaseScript.GetSwapiPlanets(context.Background())

			// then
			assert.Equal(t, cs.expectedSwapiPlanets, swapiPlanets)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FeedDatabaseScript_ParseSwapiFilmToModel(t *testing.T) {
	_, errInvalidUrl := strconv.Atoi("films")
	_, errInvalidCreated := time.Parse("2006-01-02 15:04:05", "2014-12-10")
	_, errInvalidEdited := time.Parse("2006-01-02 15:04:05", "2014-12-20")
	_, errInvalidReleaseDate := time.Parse("2006-01-02", "1977-05-25 00:00:00")

	var cases = map[string]struct {
		inputSwapiFilm model.SwapiFilm
		expectedFilm   *model.Film
		expectedErr    error
	}{
		"should return film": {
			inputSwapiFilm: model.SwapiFilm{
				Url:         "https://swapi.dev/api/films/1/",
				Created:     "2014-12-10T14:23:31.880000Z",
				Edited:      "2014-12-20T19:49:45.256000Z",
				Title:       "A New Hope",
				EpisodeID:   4,
				Director:    "George Lucas",
				ReleaseDate: "1977-05-25",
			},
			expectedFilm: &model.Film{
				ID:          1,
				CreatedAt:   time.Date(2014, 12, 10, 14, 23, 31, 0, time.UTC),
				UpdatedAt:   time.Date(2014, 12, 20, 19, 49, 45, 0, time.UTC),
				Title:       "A New Hope",
				Episode:     4,
				Director:    "George Lucas",
				ReleaseDate: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
			},
		},
		"should throw error when url is invalid": {
			inputSwapiFilm: model.SwapiFilm{Url: "https://swapi.dev/api/films/1"},
			expectedErr:    errInvalidUrl,
		},
		"should throw error when created is invalid": {
			inputSwapiFilm: model.SwapiFilm{
				Url:     "https://swapi.dev/api/films/1/",
				Created: "2014-12-10",
			},
			expectedErr: errInvalidCreated,
		},
		"should throw error when edited is invalid": {
			inputSwapiFilm: model.SwapiFilm{
				Url:     "https://swapi.dev/api/films/1/",
				Created: "2014-12-10T14:23:31.880000Z",
				Edited:  "2014-12-20",
			},
			expectedErr: errInvalidEdited,
		},
		"should throw error when release date is invalid": {
			inputSwapiFilm: model.SwapiFilm{
				Url:         "https://swapi.dev/api/films/1/",
				Created:     "2014-12-10T14:23:31.880000Z",
				Edited:      "2014-12-20T19:49:45.256000Z",
				ReleaseDate: "1977-05-25T00:00:00.000000Z",
			},
			expectedErr: errInvalidReleaseDate,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			feedDatabaseScript := &script.IFeedDatabaseScript{}

			// when
			film, err := feedDatabaseScript.ParseSwapiFilmToModel(cs.inputSwapiFilm)

			// then
			assert.Equal(t, cs.expectedFilm, film)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FeedDatabaseScript_ParseSwapiPlanetToModel(t *testing.T) {
	_, errInvalidUrl := strconv.Atoi("planets")
	_, errInvalidCreated := time.Parse("2006-01-02 15:04:05", "2014-12-09")
	_, errInvalidEdited := time.Parse("2006-01-02 15:04:05", "2014-12-20")
	climates, _ := json.Marshal([]string{"arid"})
	terrains, _ := json.Marshal([]string{"desert"})

	var cases = map[string]struct {
		inputSwapiPlanet model.SwapiPlanet
		expectedPlanet   *model.Planet
		expectedErr      error
	}{
		"should return planet": {
			inputSwapiPlanet: model.SwapiPlanet{
				Url:     "https://swapi.dev/api/planets/1/",
				Created: "2014-12-09T13:50:49.641000Z",
				Edited:  "2014-12-20T20:58:18.411000Z",
				Name:    "Tatooine",
				Climate: "arid",
				Terrain: "desert",
				Films: []string{
					"https://swapi.dev/api/films/1/",
					"https://swapi.dev/api/films/3/",
					"https://swapi.dev/api/films/4/",
					"https://swapi.dev/api/films/5/",
					"https://swapi.dev/api/films/6/",
				},
			},
			expectedPlanet: &model.Planet{
				ID:        1,
				CreatedAt: time.Date(2014, 12, 9, 13, 50, 49, 0, time.UTC),
				UpdatedAt: time.Date(2014, 12, 20, 20, 58, 18, 0, time.UTC),
				Name:      "Tatooine",
				Climates:  climates,
				Terrains:  terrains,
			},
		},
		"should throw error when url is invalid": {
			inputSwapiPlanet: model.SwapiPlanet{Url: "https://swapi.dev/api/planets/1"},
			expectedErr:      errInvalidUrl,
		},
		"should throw error when created is invalid": {
			inputSwapiPlanet: model.SwapiPlanet{
				Url:     "https://swapi.dev/api/planets/1/",
				Created: "2014-12-09",
			},
			expectedErr: errInvalidCreated,
		},
		"should throw error when edited is invalid": {
			inputSwapiPlanet: model.SwapiPlanet{
				Url:     "https://swapi.dev/api/planets/1/",
				Created: "2014-12-09T13:50:49.641000Z",
				Edited:  "2014-12-20",
			},
			expectedErr: errInvalidEdited,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			feedDatabaseScript := &script.IFeedDatabaseScript{}

			// when
			planet, err := feedDatabaseScript.ParseSwapiPlanetToModel(cs.inputSwapiPlanet)

			// then
			assert.Equal(t, cs.expectedPlanet, planet)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FeedDatabaseScript_GetIDFromUrl(t *testing.T) {
	_, errPlanets := strconv.Atoi("planets")
	_, errEmpty := strconv.Atoi("")

	var cases = map[string]struct {
		inputUrl    string
		expectedID  int
		expectedErr error
	}{
		"should return id": {
			inputUrl:   "https://swapi.dev/api/planets/1/",
			expectedID: 1,
		},
		"should throw error when id is not found": {
			inputUrl:    "https://swapi.dev/api/planets/",
			expectedID:  0,
			expectedErr: errPlanets,
		},
		"should throw error when url is empty": {
			inputUrl:    "",
			expectedID:  0,
			expectedErr: errEmpty,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			feedDatabaseScript := &script.IFeedDatabaseScript{}

			// when
			id, err := feedDatabaseScript.GetIDFromUrl(cs.inputUrl)

			// then
			assert.Equal(t, cs.expectedID, id)
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_FeedDatabaseScript_FormatStringDate(t *testing.T) {
	var cases = map[string]struct {
		inputStrDate         string
		expectedFormatedDate string
	}{
		"should return formated date": {
			inputStrDate:         "2014-12-10T11:54:13.921000Z",
			expectedFormatedDate: "2014-12-10 11:54:13",
		},
		"should return empty string when input is empty": {
			inputStrDate:         "",
			expectedFormatedDate: "",
		},
		"should return formated date when date already is formated": {
			inputStrDate:         "2014-12-10 11:54:13",
			expectedFormatedDate: "2014-12-10 11:54:13",
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			feedDatabaseScript := &script.IFeedDatabaseScript{}

			// when
			formatedDate := feedDatabaseScript.FormatStringDate(cs.inputStrDate)

			// then
			assert.Equal(t, cs.expectedFormatedDate, formatedDate)
		})
	}
}

func Test_FeedDatabaseScript_ParseToStrArrayJSON(t *testing.T) {
	var cases = map[string]struct {
		inputValue     string
		expectedValues []string
	}{
		"should return string array": {
			inputValue:     "A, B, C",
			expectedValues: []string{"A", "B", "C"},
		},
		"should return empty string array when value is empty": {
			inputValue:     "",
			expectedValues: nil,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			feedDatabaseScript := &script.IFeedDatabaseScript{}

			// when
			values := feedDatabaseScript.ParseToStrArrayJSON(cs.inputValue)

			// then
			var strValues []string
			json.Unmarshal(values, &strValues)

			assert.Equal(t, cs.expectedValues, strValues)
		})
	}
}
