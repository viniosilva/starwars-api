package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/starwars-api/internal/dto"
	"github.com/viniosilva/starwars-api/internal/exception"
	"github.com/viniosilva/starwars-api/internal/model"
	"github.com/viniosilva/starwars-api/internal/service"
)

type IPlanetController struct {
	PlanetService service.PlanetService
	Host          string
}

func (impl *IPlanetController) Configure(router *gin.RouterGroup) {
	router.GET("/planets", impl.FindPlanetsAndTotal)
	router.GET("/planets/:planetID", impl.FindPlanetByID)
	router.DELETE("/planets/:planetID", impl.DeletePlanet)
}

// @Summary find planets
// @Schemes
// @Tags planet
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param size query int false "size"
// @Param loadFilms query bool false "loadFilms"
// @Param name query string false "name"
// @Success 200 {object} dto.PlanetsResponse
// @Failure 500 {object} dto.ApiError
// @Router /api/planets [get]
func (impl *IPlanetController) FindPlanetsAndTotal(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil || size < 1 {
		size = 10
	}
	loadFilms := false
	if ctx.Query("loadFilms") == "true" {
		loadFilms = true
	}

	res, err := impl.PlanetService.FindPlanetsAndTotal(ctx, page, size, loadFilms)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	data := make([]dto.PlanetDto, res.Count)
	for i := 0; i < len(data); i += 1 {
		p := res.Data[i]

		var films []dto.FilmDto
		if p.R != nil && len(p.R.Films) > 0 {
			films = make([]dto.FilmDto, len(p.R.Films))
			for j := 0; j < len(p.R.Films); j += 1 {
				films[j] = impl.ParseFilmDto(p.R.Films[j])
			}
		}

		data[i] = impl.ParsePlanetDto(res.Data[i])
		data[i].Films = films
	}

	paramSize := ""
	if size != 10 {
		paramSize = fmt.Sprintf("&size=%d", size)
	}

	previous := ""
	if page > 1 {
		previous = fmt.Sprintf("%s?page=%d%s", impl.Host, page-1, paramSize)
	}

	next := ""
	if res.Next {
		next = fmt.Sprintf("%s?page=%d%s", impl.Host, page+1, paramSize)
	}

	ctx.JSON(http.StatusOK, dto.PlanetsResponse{
		Pagination: dto.Pagination{
			Count:    len(data),
			Total:    res.Total,
			Previous: previous,
			Next:     next,
		},
		Data: data,
	})
}

// @Summary find planet by id
// @Schemes
// @Tags planet
// @Accept json
// @Produce json
// @Param planetID path int true "Planet ID"
// @Param loadFilms query bool false "loadFilms"
// @Success 200 {object} dto.PlanetResponse
// @Failure 400 {object} dto.ApiError
// @Failure 404 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Router /api/planets/{planetID} [get]
func (impl *IPlanetController) FindPlanetByID(ctx *gin.Context) {
	planetID, err := strconv.Atoi(ctx.Param("planetID"))
	if err != nil || planetID < 1 {
		ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: "invalid planet id"})
		return
	}
	loadFilms := false
	if ctx.Query("loadFilms") == "true" {
		loadFilms = true
	}

	planet, err := impl.PlanetService.FindPlanetByID(ctx, planetID, loadFilms)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			ctx.JSON(http.StatusNotFound, dto.ApiError{Error: fmt.Sprintf("planet %d not found", planetID)})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	data := impl.ParsePlanetDto(planet)

	var films []dto.FilmDto
	if planet.R != nil && len(planet.R.Films) > 0 {
		films = make([]dto.FilmDto, len(planet.R.Films))
		for j := 0; j < len(planet.R.Films); j += 1 {
			films[j] = impl.ParseFilmDto(planet.R.Films[j])
		}
		data.Films = films
	}

	ctx.JSON(http.StatusOK, dto.PlanetResponse{Data: data})
}

// @Summary delete planet
// @Schemes
// @Tags planet
// @Accept json
// @Produce json
// @Param planetID path int true "Planet ID"
// @Success 204 ""
// @Failure 400 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Router /api/planets/{planetID} [delete]
func (impl *IPlanetController) DeletePlanet(ctx *gin.Context) {
	planetID, err := strconv.Atoi(ctx.Param("planetID"))
	if err != nil || planetID < 1 {
		ctx.JSON(http.StatusBadRequest, dto.ApiError{Error: "invalid planet id"})
		return
	}

	err = impl.PlanetService.DeletePlanet(ctx, planetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (impl *IPlanetController) ParsePlanetDto(planet *model.Planet) dto.PlanetDto {
	var climates []string
	json.Unmarshal(planet.Climates, &climates)

	var terrains []string
	json.Unmarshal(planet.Terrains, &terrains)

	return dto.PlanetDto{
		ID:        planet.ID,
		CreatedAt: planet.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: planet.UpdatedAt.Format("2006-01-02 15:04:05"),
		Name:      planet.Name,
		Climates:  climates,
		Terrains:  terrains,
	}
}

func (impl *IPlanetController) ParseFilmDto(film *model.Film) dto.FilmDto {
	return dto.FilmDto{
		ID:          film.ID,
		CreatedAt:   film.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   film.UpdatedAt.Format("2006-01-02 15:04:05"),
		Title:       film.Title,
		Episode:     int(film.Episode),
		Director:    film.Director,
		ReleaseDate: film.ReleaseDate.Format("2006-01-02"),
	}
}
