package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viniosilva/starwars-api/docs"
	"github.com/viniosilva/starwars-api/internal/config"
	"github.com/viniosilva/starwars-api/internal/controller"
	"github.com/viniosilva/starwars-api/internal/service"
)

// @title		Star Wars API
// @version		1.0
// @BasePath	/api
func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename: "log/logrus.log",
		MaxSize:  50, // megabytes
	}))

	c := config.LoadConfig()

	r := gin.Default()
	r.Use(config.GinLogger())

	router := r.Group("/api")

	healthService := &service.IHealthService{}
	planetService := &service.IPlanetService{}

	healthController := &controller.IHealthController{HealthService: healthService}
	planetController := &controller.IPlanetController{PlanetService: planetService}

	healthController.Configure(router)
	planetController.Configure(router)

	host := fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
	docs.SwaggerInfo.Host = host
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logrus.WithFields(logrus.Fields{"trace": "main"}).Infof("listening on %s", host)
	r.Run(host)
}
