package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ovidiuz/device_manager/query_service/gateways/interfaces"
	"github.com/sirupsen/logrus"
)

type MetricsHandler struct {
	metricsRepo interfaces.MetricsRepo
}

func NewMetricsHandler(metricsRepo interfaces.MetricsRepo) *MetricsHandler {
	return &MetricsHandler{metricsRepo: metricsRepo}
}

func (m *MetricsHandler) RegisterRoutes(ws *fiber.App) {
	ws.Get("/metrics", m.GetAllMeasurements)
	ws.Get("/metrics/:id", m.GetAllMeasurementsForSensor)
}

func (m *MetricsHandler) GetAllMeasurements(ctx *fiber.Ctx) (err error) {
	var apiErr *fiber.Error
	defer func() {
		if apiErr != nil {
			ctx.Status(apiErr.Code)
			err = ctx.JSON(apiErr)
		}
	}()

	measurements, err := m.metricsRepo.GetAllMeasurements(ctx.UserContext())
	if err != nil {
		logrus.WithError(err).Error("could not get measurements")
		apiErr = fiber.ErrInternalServerError
		return
	}

	return ctx.JSON(measurements)
}

func (m *MetricsHandler) GetAllMeasurementsForSensor(ctx *fiber.Ctx) (err error) {
	var apiErr *fiber.Error
	sensorID := ctx.Params("id")

	defer func() {
		if apiErr != nil {
			ctx.Status(apiErr.Code)
			err = ctx.JSON(apiErr)
		}
	}()

	measurements, err := m.metricsRepo.GetMeasurementForSensor(ctx.UserContext(), sensorID)
	if err != nil {
		logrus.WithError(err).Error("could not get measurements")
		apiErr = fiber.ErrInternalServerError
		return
	}

	return ctx.JSON(measurements)
}
