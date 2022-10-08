package service

import (
	"context"
)

//go:generate mockgen -destination=../../mock/health_service_mock.go -package=mock . HealthService
type HealthService interface {
	Ping(ctx context.Context) error
}

type IHealthService struct{}

func (impl *IHealthService) Ping(ctx context.Context) error {
	// TODO
	return nil
}
