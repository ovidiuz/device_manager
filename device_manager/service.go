package main

import (
	"fmt"
	"net/http"

	"github.com/ovidiuz/device_manager/gateways/repositories"
	"github.com/ovidiuz/device_manager/usecases"

	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo/v4"
	"github.com/ovidiuz/device_manager/domain"
	"github.com/ovidiuz/device_manager/gateways/api"
	"github.com/ovidiuz/device_manager/gateways/interfaces"

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

	// init Cabin RBAC Policy Management
	casbinEnforcer := initCasbin(conf, sqlDB)

	// init managers
	userManager := usecases.NewUserManager(userSQLRepo, casbinEnforcer, conf.JwtTTL)

	// init API route handlers
	authHandler := api.NewAuthHandler(userManager, conf.JwtTTL, conf.HTTPSEnabled)
	userHandler := api.NewUserHandler(userManager, authHandler.GetAuthMiddleware())
	apiRouteHandlers := []interfaces.RouteHandler{
		authHandler,
		userHandler,
	}

	// setup & start the HTTP server
	startHTTPServer(conf.ServicePort, apiRouteHandlers)
}

func startHTTPServer(port int, routeHandlers []interfaces.RouteHandler) {
	ws := echo.New()
	ws.HTTPErrorHandler = customHTTPErrorHandler(ws)
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

func initCasbin(conf *domain.ServiceConfig, db *sqlx.DB) casbin.IEnforcer {
	//adapter, err := sqlxadapter.NewAdapter(db, conf.CasbinTableName)
	//if err != nil {
	//	logrus.WithError(err).Fatal("could not create Casbin Sqlx adapter")
	//}
	//
	//enforcer, err := casbin.NewEnforcer(conf.CasbinModelFile, adapter)
	//if err != nil {
	//	logrus.WithError(err).Fatal("could not create Casbin enforcer")
	//}
	//
	//return enforcer
	return nil
}

func customHTTPErrorHandler(ws *echo.Echo) func(err error, ctx echo.Context) {
	return func(err error, ctx echo.Context) {
		defer func() {
			ws.DefaultHTTPErrorHandler(err, ctx)
		}()

		if err == nil {
			return
		}
		_, ok := err.(*echo.HTTPError)
		if ok {
			return
		}

		switch err {
		case domain.ErrNotFound:
			err = echo.ErrNotFound
		case domain.ErrUnauthorized:
			err = echo.ErrUnauthorized
		default:
			err = echo.ErrInternalServerError
		}
	}
}
