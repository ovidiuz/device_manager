package repositories

import (
	"context"
	"fmt"

	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type MetricsInfluxRepo struct {
	client   influx.Client
	queryApi api.QueryAPI
}

func NewMetricsInfluxRepo(client influx.Client, queryApi api.QueryAPI) *MetricsInfluxRepo {
	return &MetricsInfluxRepo{client: client, queryApi: queryApi}
}

func (r *MetricsInfluxRepo) GetAllMeasurements(ctx context.Context) (interface{}, error) {
	rows := make([]string, 0)
	queryStmt := `from(bucket: "my-bucket")
	  |> range(start: -48h, stop: now())
	  |> filter(fn: (r) => r["_measurement"] == "mqtt_consumer")`

	result, err := r.queryApi.Query(ctx, queryStmt)
	if err == nil {
		for result.Next() {
			rows = append(rows, result.Record().String())
		}
		if result.Err() != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return rows, nil
}

func (r *MetricsInfluxRepo) GetMeasurementForSensor(ctx context.Context, sensor string) (interface{}, error) {
	rows := make([]string, 0)
	queryFormat := `from(bucket: "my-bucket")
	  |> range(start: -24h, stop: now())
	  |> filter(fn: (r) => r["_measurement"] == "mqtt_consumer")
	  |> filter(fn: (r) => r["topic"] == "%s")`
	queryStmt := fmt.Sprintf(queryFormat, sensor)

	result, err := r.queryApi.Query(ctx, queryStmt)
	if err == nil {
		for result.Next() {
			rows = append(rows, result.Record().String())
		}
		if result.Err() != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return rows, nil
}
