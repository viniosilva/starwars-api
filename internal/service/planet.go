package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/starwars-api/internal/model"
)

//go:generate mockgen -destination=../../mock/planet_service_mock.go -package=mock . PlanetService
type PlanetService interface {
	CreatePlanets(ctx context.Context, planets []*model.Planet) error
	CreateRelationshipFilmsToPlanets(ctx context.Context, relationships map[int][]int) error
	FindPlanets(ctx context.Context, page, size int) ([]model.Planet, int, error)
	FindPlanetByID(ctx context.Context, planetID int) (*model.Planet, error)
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

func (impl *IPlanetService) FindPlanets(ctx context.Context, page, size int) ([]model.Planet, int, error) {
	// TODO
	return nil, 0, nil
}

func (impl *IPlanetService) FindPlanetByID(ctx context.Context, planetID int) (*model.Planet, error) {
	// TODO
	return nil, nil
}

func (impl *IPlanetService) DeletePlanet(ctx context.Context, planetID int) error {
	//TODO
	return nil
}
