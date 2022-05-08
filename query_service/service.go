package main

import (
	"context"
	"fmt"

	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/ovidiuz/query_service/domain"
	"github.com/ovidiuz/query_service/gateways/interfaces"
	"github.com/ovidiuz/query_service/gateways/repositories"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Launching the device manager service")

	// Load the configurations
	conf, err := domain.InitConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	metricsRepo := initMetricsRepo(conf)

	ctx := context.Background()
	results, err := metricsRepo.GetAllMeasurements(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(results)
}

func initMetricsRepo(conf *domain.ServiceConfig) interfaces.MetricsRepo {
	client := influx.NewClient(conf.InfluxHost, conf.InfluxToken)
	queryApi := client.QueryAPI(conf.InfluxOrg)
	return repositories.NewMetricsInfluxRepo(client, queryApi)
}
