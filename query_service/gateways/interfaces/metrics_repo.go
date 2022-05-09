package interfaces

import (
	"context"

	"github.com/ovidiuz/device_manager/query_service/domain"
)

type MetricsRepo interface {
	GetAllMeasurements(ctx context.Context) ([]domain.Measurement, error)
	GetMeasurementForSensor(ctx context.Context, sensor string) ([]domain.Measurement, error)
}
