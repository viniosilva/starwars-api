package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/viniosilva/starwars-api/internal/model"
)

//go:generate mockgen -destination=../../mock/film_service_mock.go -package=mock . FilmService
type FilmService interface {
	CreateFilms(ctx context.Context, Films []*model.Film) error
}

type IFilmService struct {
	DB *sql.DB
}

func (impl *IFilmService) CreateFilms(ctx context.Context, films []*model.Film) error {
	values := make([]string, len(films))
	for i := 0; i < len(values); i += 1 {
		f := films[i]
		values[i] = fmt.Sprintf("(%d, '%s', '%s', '%s', '%s', %d, '%s')",
			f.ID,
			f.CreatedAt.Format("2006-01-02 15:04:05"),
			f.UpdatedAt.Format("2006-01-02 15:04:05"),
			f.Director,
			f.Title,
			f.Episode,
			f.ReleaseDate.Format("2006-01-02"),
		)
	}

	columns := []string{
		model.FilmColumns.ID,
		model.FilmColumns.CreatedAt,
		model.FilmColumns.UpdatedAt,
		model.FilmColumns.Director,
		model.FilmColumns.Title,
		model.FilmColumns.Episode,
		model.FilmColumns.ReleaseDate,
	}

	query := fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES %s;",
		model.TableNames.Films, strings.Join(columns, ", "), strings.Join(values, ",\n"))

	_, err := impl.DB.ExecContext(ctx, query)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "internal.service.film.create_films:tx.exec_context"}).Error(err)
		return err
	}

	return nil
}
