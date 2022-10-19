package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/viniosilva/starwars-api/docs"
	"github.com/viniosilva/starwars-api/internal/config"
	"github.com/viniosilva/starwars-api/internal/controller"
	"github.com/viniosilva/starwars-api/internal/request"
	"github.com/viniosilva/starwars-api/internal/script"
	"github.com/viniosilva/starwars-api/internal/service"
)

const (
	LOGS_PATH         = "log/logrus.log"
	ARG_FEED_DATABASE = "feed_database"
)

func main() {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename: LOGS_PATH,
		MaxSize:  50, // megabytes
	}))

	c := config.LoadConfig()
	host := fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)

	db_conn_string := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
		c.MySQL.Username, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.Database)
	db, err := sql.Open("mysql", db_conn_string)
	if err != nil {
		panic(err)
	}

	healthService := &service.IHealthService{DB: db}
	filmService := &service.IFilmService{DB: db}
	planetService := &service.IPlanetService{DB: db}

	if len(os.Args) > 1 && os.Args[1] == ARG_FEED_DATABASE {
		go runScript(filmService, planetService)
	} else {
		go runApi(host, healthService, planetService)
	}

	<-gracefulShutdown
	if err := db.Close(); err != nil {
		logrus.WithField("trace", "main").Error(err)
	}

	logrus.WithField("trace", "main").Info("shutdown")
}

func runScript(filmService service.FilmService, planetService service.PlanetService) {
	swapi := &request.ISwapiRequest{}
	feedDatabase := &script.IFeedDatabaseScript{
		Swapi:         swapi,
		FilmService:   filmService,
		PlanetService: planetService,
	}

	if err := feedDatabase.Execute(); err != nil {
		panic(err)
	}
	os.Exit(1)
}

// @title		Star Wars API
// @version		1.0
// @BasePath	/api
func runApi(host string, healthService service.HealthService, planetService service.PlanetService) {
	r := gin.Default()
	r.Use(config.GinLogger())

	router := r.Group("/api")

	healthController := &controller.IHealthController{HealthService: healthService}
	planetController := &controller.IPlanetController{
		Host:          fmt.Sprintf("http://%s/api/planets", host),
		PlanetService: planetService,
	}

	healthController.Configure(router)
	planetController.Configure(router)

	docs.SwaggerInfo.Host = host
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logrus.WithFields(logrus.Fields{"trace": "main"}).Infof("listening on %s", host)
	r.Run(host)
}
