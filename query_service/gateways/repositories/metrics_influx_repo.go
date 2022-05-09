package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/ovidiuz/device_manager/query_service/domain"
)

type MetricsInfluxRepo struct {
	client   influx.Client
	queryApi api.QueryAPI
}

func NewMetricsInfluxRepo(client influx.Client, queryApi api.QueryAPI) *MetricsInfluxRepo {
	return &MetricsInfluxRepo{client: client, queryApi: queryApi}
}

func (r *MetricsInfluxRepo) GetAllMeasurements(ctx context.Context) ([]domain.Measurement, error) {
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
			logrus.WithError(result.Err()).Error("could not get measurements")
			return nil, err
		}
	} else {
		logrus.WithError(err).Error("could not get measurements")
		return nil, err
	}

	return convertRowsToMeasurements(rows), nil
}

func (r *MetricsInfluxRepo) GetMeasurementForSensor(ctx context.Context, sensor string) ([]domain.Measurement, error) {
	rows := make([]string, 0)
	queryFormat := `from(bucket: "my-bucket")
	  |> range(start: -48h, stop: now())
	  |> filter(fn: (r) => r["_measurement"] == "mqtt_consumer")
	  |> filter(fn: (r) => r["topic"] == "sensors/%s")`
	queryStmt := fmt.Sprintf(queryFormat, sensor)

	result, err := r.queryApi.Query(ctx, queryStmt)
	if err == nil {
		for result.Next() {
			rows = append(rows, result.Record().String())
		}
		if result.Err() != nil {
			logrus.WithError(result.Err()).Error("could not get measurements")
			return nil, err
		}
	} else {
		logrus.WithError(err).Error("could not get measurements")
		return nil, err
	}

	return convertRowsToMeasurements(rows), nil
}

func convertRowsToMeasurements(rows []string) []domain.Measurement {
	measurements := make([]domain.Measurement, 0)
	for _, row := range rows {
		measurement := convertRowToMeasurement(row)
		measurements = append(measurements, measurement)
	}

	return measurements
}

func convertRowToMeasurement(row string) domain.Measurement {
	rowParts := strings.Split(row, ",")

	fmt.Println(row)
	fmt.Println(rowParts)
	parts := strings.Split(rowParts[0], ":")
	field := parts[1]

	parts = strings.Split(rowParts[1], ":")
	measurement := parts[1]

	parts = strings.SplitN(rowParts[2], ":", 2)
	start := parts[1]

	parts = strings.SplitN(rowParts[3], ":", 2)
	stop := parts[1]

	parts = strings.SplitN(rowParts[4], ":", 2)
	timeV := parts[1]

	parts = strings.Split(rowParts[5], ":")
	value, _ := strconv.Atoi(parts[1])

	parts = strings.Split(rowParts[6], ":")
	host := parts[1]

	parts = strings.Split(rowParts[7], ":")
	result := parts[1]

	parts = strings.Split(rowParts[8], ":")
	table, _ := strconv.Atoi(parts[1])

	parts = strings.Split(rowParts[9], ":")
	topic := parts[1]

	return domain.Measurement{
		Field:       field,
		Measurement: measurement,
		Start:       start,
		Stop:        stop,
		Time:        timeV,
		Value:       value,
		Host:        host,
		Result:      result,
		Table:       table,
		Topic:       topic,
	}
}
