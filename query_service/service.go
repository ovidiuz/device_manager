package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/ovidiuz/device_manager/query_service/domain"
	"github.com/ovidiuz/device_manager/query_service/gateways/api"
	"github.com/ovidiuz/device_manager/query_service/gateways/interfaces"
	"github.com/ovidiuz/device_manager/query_service/gateways/repositories"
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
	metricsHandler := api.NewMetricsHandler(metricsRepo)

	ws := fiber.New()
	metricsHandler.RegisterRoutes(ws)

	listenAddr := fmt.Sprintf(":%d", conf.WebServicePort)
	if err := ws.Listen(listenAddr); err != nil {
		panic(err)
	}
}

func initMetricsRepo(conf *domain.ServiceConfig) interfaces.MetricsRepo {
	client := influx.NewClient(conf.InfluxHost, conf.InfluxToken)
	queryApi := client.QueryAPI(conf.InfluxOrg)
	return repositories.NewMetricsInfluxRepo(client, queryApi)
}
