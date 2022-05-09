package interfaces

import (
	fiber "github.com/gofiber/fiber/v2"
)

type RouteHandler interface {
	RegisterRoutes(ws *fiber.App)
}
