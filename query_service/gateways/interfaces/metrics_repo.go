package interfaces

import "context"

type MetricsRepo interface {
	GetAllMeasurements(ctx context.Context) (interface{}, error)
	GetMeasurementForSensor(ctx context.Context, sensor string) (interface{}, error)
}
