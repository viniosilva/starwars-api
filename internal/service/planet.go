package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/starwars-api/internal/dto"
	"github.com/viniosilva/starwars-api/internal/exception"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

//go:generate mockgen -destination=../../mock/planet_service_mock.go -package=mock . PlanetService
type PlanetService interface {
	CreatePlanets(ctx context.Context, planets []*model.Planet) error
	CreateRelationshipFilmsToPlanets(ctx context.Context, relationships map[int][]int) error
	FindPlanetsAndTotal(ctx context.Context, page, size int, loadFilms bool, opts ...Option) (dto.FindPlanetsAndTotalResult, error)
	FindPlanetByID(ctx context.Context, planetID int, loadFilms bool) (*model.Planet, error)
	DeletePlanet(ctx context.Context, planetID int) error
}

type IPlanetService struct {
	DB *sql.DB
}

func (impl *IPlanetService) CreatePlanets(ctx context.Context, planets []*model.Planet) error {
	values := make([]string, len(planets))
	for i := 0; i < len(values); i += 1 {
		p := planets[i]
		values[i] = fmt.Sprintf("(%d, '%s', '%s', '%s', '%s', '%s')",
			p.ID,
			p.CreatedAt.Format("2006-01-02 15:04:05"),
			p.UpdatedAt.Format("2006-01-02 15:04:05"),
			p.Name,
			p.Climates.String(),
			p.Terrains.String(),
		)
	}

	columns := []string{
		model.PlanetColumns.ID,
		model.PlanetColumns.CreatedAt,
		model.PlanetColumns.UpdatedAt,
		model.PlanetColumns.Name,
		model.PlanetColumns.Climates,
		model.PlanetColumns.Terrains,
	}
	query := fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES %s;",
		model.TableNames.Planets, strings.Join(columns, ", "), strings.Join(values, ",\n"))

	_, err := impl.DB.ExecContext(ctx, query)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.create_planets:tx.exec_context"}).Error(err)
		return err
	}

	return nil
}

func (impl *IPlanetService) CreateRelationshipFilmsToPlanets(ctx context.Context, relationships map[int][]int) error {
	values := []string{}
	for planetID, filmIDs := range relationships {
		for _, filmID := range filmIDs {
			values = append(values, fmt.Sprintf("(%d, %d)", planetID, filmID))
		}
	}

	query := fmt.Sprintf("INSERT IGNORE INTO %s (planet_id, film_id) VALUES %s;",
		model.TableNames.PlanetsFilms, strings.Join(values, ",\n"))

	_, err := impl.DB.ExecContext(ctx, query)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.create_relationship_films_to_planets:db.exec_context"}).Error(err)

		return err
	}

	return nil
}

func (impl *IPlanetService) FindPlanetsAndTotal(ctx context.Context, page, size int, loadFilms bool, opts ...Option) (dto.FindPlanetsAndTotalResult, error) {
	offset := 0
	if page > 1 {
		offset = size * (page - 1)
	}

	tx, err := impl.DB.BeginTx(ctx, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.find_planets_and_total:db.begin_tx"}).Error(err)
		return dto.FindPlanetsAndTotalResult{}, err
	}

	whereIsNotDeleted := qm.Where(fmt.Sprintf("%s IS NULL", model.PlanetColumns.DeletedAt))
	qms := []qm.QueryMod{
		qm.Limit(size + 1),
		qm.Offset(offset),
		whereIsNotDeleted,
	}

	where, arg := GetOptionWhere(opts)
	if where != "" {
		qms = append(qms, qm.And(where, arg))
	}

	if loadFilms {
		qms = append(qms, qm.Load("Films"))
	}

	planets, err := model.Planets(qms...).All(ctx, tx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.FindPlanets:planets.all"}).Error(err)
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.FindPlanets:tx.rollback"}).Error(err)
			return dto.FindPlanetsAndTotalResult{}, err
		}

		return dto.FindPlanetsAndTotalResult{}, err
	}

	qms = []qm.QueryMod{whereIsNotDeleted}
	if where != "" {
		qms = append(qms, qm.And(where, arg))
	}

	total, err := model.Planets(qms...).Count(ctx, tx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.FindPlanets:planets.count"}).Error(err)
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.FindPlanets:tx.rollback"}).Error(err)
			return dto.FindPlanetsAndTotalResult{}, err
		}

		return dto.FindPlanetsAndTotalResult{}, err
	}

	if err := tx.Commit(); err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.FindPlanets:tx.commit"}).Error(err)
		return dto.FindPlanetsAndTotalResult{}, err
	}

	data := planets
	next := false
	if len(planets) > size {
		next = true
		data = planets[:size]
	}

	return dto.FindPlanetsAndTotalResult{
		Total: total,
		Count: len(data),
		Next:  next,
		Data:  data,
	}, nil
}

func (impl *IPlanetService) FindPlanetByID(ctx context.Context, planetID int, loadFilms bool) (*model.Planet, error) {
	qms := []qm.QueryMod{
		qm.Where(fmt.Sprintf("%s = ?", model.PlanetColumns.ID), planetID),
		qm.Where(fmt.Sprintf("%s IS NULL", model.PlanetColumns.DeletedAt)),
	}
	if loadFilms {
		qms = append(qms, qm.Load("Films"))
	}

	planet, err := model.Planets(qms...).One(ctx, impl.DB)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, &exception.NotFoundException{
				Message: fmt.Sprintf("planet %d not found", planetID),
			}
		}

		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.find_planet_by_id:planets.one"}).Error(err)
		return nil, err
	}
	return planet, nil
}

func (impl *IPlanetService) DeletePlanet(ctx context.Context, planetID int) error {
	_, err := model.Planets(qm.Where(fmt.Sprintf("%s = ?", model.PlanetColumns.ID), planetID)).
		UpdateAll(ctx, impl.DB, model.M{model.PlanetColumns.DeletedAt: time.Now()})
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.planet.delete_planet:planets.update_all"}).Error(err)
		return err
	}

	return nil
}
