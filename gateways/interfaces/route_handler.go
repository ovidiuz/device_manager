package interfaces

import "github.com/labstack/echo"

type RouteHandler interface {
	RegisterRoutes(ws *echo.Echo)
}
