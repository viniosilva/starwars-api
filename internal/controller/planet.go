package controller

import (
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
}

func (impl *IPlanetController) Configure(router *gin.RouterGroup) {
	router.GET("/planets", impl.FindPlanets)
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
// @Success 200 {object} dto.PlanetsResponse
// @Failure 500 {object} dto.ApiError
// @Router /api/planets [get]
func (impl *IPlanetController) FindPlanets(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil || size < 1 {
		size = 10
	}

	planets, total, err := impl.PlanetService.FindPlanets(ctx, page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	data := make([]dto.PlanetDto, len(planets))
	for i := 0; i < len(planets); i += 1 {
		// films := make([]dto.FilmDto, len(planets[i].Films))
		// for j := 0; j < len(planets[i].Films); j += 1 {
		// 	films[j] = impl.ParseFilmDto(&planets[i].Films[j])
		// }

		data[i] = impl.ParsePlanetDto(&planets[i])
		// data[i].Films = films
	}

	ctx.JSON(http.StatusOK, dto.PlanetsResponse{
		Pagination: dto.Pagination{
			Count: len(data),
			Total: total,
		},
		Data: data})
}

// @Summary find planet by id
// @Schemes
// @Tags planet
// @Accept json
// @Produce json
// @Param planetID path int true "Planet ID"
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

	planet, err := impl.PlanetService.FindPlanetByID(ctx, planetID)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); ok {
			ctx.JSON(http.StatusNotFound, dto.ApiError{Error: fmt.Sprintf("planet %d not found", planetID)})
		}
		ctx.JSON(http.StatusInternalServerError, dto.ApiError{Error: "internal server error"})
		return
	}

	// films := make([]dto.FilmDto, len(planet.Films))
	// for j := 0; j < len(planet.Films); j += 1 {
	// 	films[j] = impl.ParseFilmDto(&planet.Films[j])
	// }

	data := impl.ParsePlanetDto(planet)
	// data.Films = films

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
	return dto.PlanetDto{
		ID:        planet.ID,
		CreatedAt: planet.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: planet.CreatedAt.Format("2006-01-02 15:04:05"),
		Name:      planet.Name,
		// Climates:  planet.Climates,
		// Terrains:  planet.Terrains,
	}
}

func (impl *IPlanetController) ParseFilmDto(film *model.Film) dto.FilmDto {
	return dto.FilmDto{
		ID:          film.ID,
		CreatedAt:   film.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   film.UpdatedAt.Format("2006-01-02 15:04:05"),
		Title:       film.Title,
		Director:    film.Director,
		ReleaseDate: film.ReleaseDate.Format("2006-01-02"),
	}
}
