package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ovidiuz/device_manager/domain"
	"github.com/ovidiuz/device_manager/gateways/api"
	"github.com/ovidiuz/device_manager/gateways/interfaces"
	"github.com/ovidiuz/device_manager/gateways/repositories"
	"github.com/ovidiuz/device_manager/usecases"

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

	sqlDB, err := initSQL(conf)
	if err != nil {
		logrus.WithError(err).Fatal("could not initialize an SQL connection")
	}

	// init repos
	userSQLRepo := repositories.NewUserSQLRepo(sqlDB)

	// init managers
	userManager := usecases.NewUserManager(userSQLRepo)

	// init API route handlers
	apiRouteHandlers := []interfaces.RouteHandler{
		api.NewAuthHandler(userManager),
		api.NewUserHandler(userManager),
	}

	// setup & start the HTTP server
	startHTTPServer(conf.ServicePort, apiRouteHandlers)
}

func startHTTPServer(port int, routeHandlers []interfaces.RouteHandler) {
	ws := echo.New()
	for _, handler := range routeHandlers {
		handler.RegisterRoutes(ws)
	}

	address := fmt.Sprintf(":%d", port)
	if err := ws.Start(address); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("HTTP server stopped")
	}

	// will probably very rarely reach this log statement
	logrus.Info("HTTP server gracefully shut down")
}
