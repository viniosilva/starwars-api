package service

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../../mock/health_service_mock.go -package=mock . HealthService
type HealthService interface {
	Ping(ctx context.Context) error
}

type IHealthService struct {
	DB *sql.DB
}

func (impl *IHealthService) Ping(ctx context.Context) error {
	err := impl.DB.PingContext(ctx)
	if err != nil {
		logrus.WithField("trace", "internal.service.health.ping").Error(err)
	}

	return err
}
