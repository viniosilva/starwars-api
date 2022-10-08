package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniosilva/starwars-api/internal/dto"
	"github.com/viniosilva/starwars-api/internal/service"
)

type IHealthController struct {
	HealthService service.HealthService
}

func (impl *IHealthController) Configure(router *gin.RouterGroup) {
	router.GET("/healthcheck", impl.Ping)
}

// @Summary healthcheck
// @Schemes
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} dto.HealthResponse
// @Router /api/healthcheck [get]
func (impl *IHealthController) Ping(ctx *gin.Context) {
	if err := impl.HealthService.Ping(ctx); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, dto.HealthResponse{Status: dto.HealshStatusDown})
		return
	}

	ctx.JSON(http.StatusOK, dto.HealthResponse{Status: dto.HealshStatusUp})
}
