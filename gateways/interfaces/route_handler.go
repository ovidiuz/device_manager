package interfaces

import "github.com/labstack/echo/v4"

type RouteHandler interface {
	RegisterRoutes(ws *echo.Echo)
}
