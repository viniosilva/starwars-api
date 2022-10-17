package script

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/request"
	"github.com/viniosilva/starwars-api/internal/service"
)

type IFeedDatabaseScript struct {
	Swapi         request.SwapiRequest
	FilmService   service.FilmService
	PlanetService service.PlanetService
}

const TRACE_EXECUTE = "internal.script.execute"

func (impl *IFeedDatabaseScript) Execute() error {
	logrus.WithFields(logrus.Fields{"trace": TRACE_EXECUTE}).Info("starting")

	ctx := context.Background()

	swapiFilms, err := impl.GetSwapiFilms(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": fmt.Sprintf("%s:get_swapi_films", TRACE_EXECUTE)}).Error(err)

		return err
	}

	films := make([]*model.Film, len(swapiFilms))
	for i := 0; i < len(films); i += 1 {
		film, err := impl.ParseSwapiFilmToModel(swapiFilms[i])
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"trace":    fmt.Sprintf("%s:parse_swapi_film_to_model", TRACE_EXECUTE),
				"film_url": swapiFilms[i].Url,
			}).Error(err)

			return err
		}

		films[i] = film
	}

	swapiPlanets, err := impl.GetSwapiPlanets(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": fmt.Sprintf("%s:get_swapi_planets", TRACE_EXECUTE)}).Error(err)

		return err
	}

	planets := make([]*model.Planet, len(swapiPlanets))
	for i := 0; i < len(planets); i += 1 {
		planet, err := impl.ParseSwapiPlanetToModel(swapiPlanets[i])
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"trace":      fmt.Sprintf("%s:parse_swapi_planet_to_model", TRACE_EXECUTE),
				"planet_url": swapiPlanets[i].Url,
			}).Error(err)

			return err
		}

		planets[i] = planet
	}

	relationships := map[int][]int{}
	for _, p := range swapiPlanets {
		planetID, err := impl.GetIDFromUrl(p.Url)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": fmt.Sprintf("%s:get_id_from_url", TRACE_EXECUTE)}).Error(err)
			return err
		}

		for _, f := range p.Films {
			filmID, err := impl.GetIDFromUrl(f)
			if err != nil {
				logrus.WithFields(logrus.Fields{"trace": fmt.Sprintf("%s:get_id_from_url", TRACE_EXECUTE)}).Error(err)

				return err
			}
			relationships[planetID] = append(relationships[planetID], filmID)
		}
	}

	err = impl.FilmService.CreateFilms(ctx, films)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": fmt.Sprintf("%s:create_films", TRACE_EXECUTE)}).Error(err)

		return err
	}

	err = impl.PlanetService.CreatePlanets(ctx, planets)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": fmt.Sprintf("%s:create_planets", TRACE_EXECUTE)}).Error(err)

		return err
	}

	err = impl.PlanetService.CreateRelationshipFilmsToPlanets(ctx, relationships)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": fmt.Sprintf("%s:create_planets", TRACE_EXECUTE)}).Error(err)

		return err
	}

	logrus.WithFields(logrus.Fields{"trace": TRACE_EXECUTE}).Info("finished")
	return nil
}

func (impl *IFeedDatabaseScript) GetSwapiFilms(ctx context.Context) ([]model.SwapiFilm, error) {
	page := 1
	films := []model.SwapiFilm{}
	for page != 0 {
		res, err := impl.Swapi.GetFilms(ctx, page)
		if err != nil {
			return nil, err
		}

		page += 1
		if res.Next == "" {
			page = 0
		}

		films = append(films, res.Results...)
	}

	return films, nil
}

func (impl *IFeedDatabaseScript) GetSwapiPlanets(ctx context.Context) ([]model.SwapiPlanet, error) {
	page := 1
	planets := []model.SwapiPlanet{}
	for page != 0 {
		res, err := impl.Swapi.GetPlanets(ctx, page)
		if err != nil {
			return nil, err
		}

		page += 1
		if res.Next == "" {
			page = 0
		}

		planets = append(planets, res.Results...)

	}

	return planets, nil
}

func (impl *IFeedDatabaseScript) ParseSwapiFilmToModel(swapiFilm model.SwapiFilm) (*model.Film, error) {
	var err error
	film := &model.Film{
		Title:    swapiFilm.Title,
		Episode:  int8(swapiFilm.EpisodeID),
		Director: swapiFilm.Director,
	}

	film.ID, err = impl.GetIDFromUrl(swapiFilm.Url)
	if err != nil {
		return nil, err
	}
	if film.CreatedAt, err = time.Parse("2006-01-02 15:04:05", impl.FormatStringDate(swapiFilm.Created)); err != nil {
		return nil, err
	}
	if film.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", impl.FormatStringDate(swapiFilm.Edited)); err != nil {
		return nil, err
	}
	if film.ReleaseDate, err = time.Parse("2006-01-02", impl.FormatStringDate(swapiFilm.ReleaseDate)); err != nil {
		return nil, err
	}

	return film, nil
}

func (impl *IFeedDatabaseScript) ParseSwapiPlanetToModel(swapiPlanet model.SwapiPlanet) (*model.Planet, error) {
	var err error

	planet := &model.Planet{
		Name:     swapiPlanet.Name,
		Climates: impl.ParseToStrArrayJSON(swapiPlanet.Climate),
		Terrains: impl.ParseToStrArrayJSON(swapiPlanet.Terrain),
	}

	planet.ID, err = impl.GetIDFromUrl(swapiPlanet.Url)
	if err != nil {
		return nil, err
	}
	if planet.CreatedAt, err = time.Parse("2006-01-02 15:04:05", impl.FormatStringDate(swapiPlanet.Created)); err != nil {
		return nil, err
	}
	if planet.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", impl.FormatStringDate(swapiPlanet.Edited)); err != nil {
		return nil, err
	}

	return planet, nil
}

func (impl *IFeedDatabaseScript) GetIDFromUrl(url string) (int, error) {
	e := strings.Split(url, "/")

	strID := ""
	if len(e) > 1 {
		strID = e[(len(e) - 2)]
	}

	id, err := strconv.Atoi(strID)
	return id, err
}

func (impl *IFeedDatabaseScript) FormatStringDate(strDate string) string {
	date := strings.Replace(strDate, "T", " ", 1)
	date = strings.Split(date, ".")[0]

	return date
}

func (impl *IFeedDatabaseScript) ParseToStrArrayJSON(value string) []byte {
	if value == "" {
		return []byte{}
	}

	values := strings.Split(value, ",")
	for i := 0; i < len(values); i += 1 {
		values[i] = strings.TrimSpace(values[i])
	}

	b, err := json.Marshal(values)
	if err != nil {
		return []byte{}
	}

	return b
}
