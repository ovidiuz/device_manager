package api

import (
	"github.com/labstack/echo/v4"
	"github.com/ovidiuz/device_manager/usecases"
)

type UserHandler struct {
	middleware  []echo.MiddlewareFunc
	userManager *usecases.UserManager
}

func NewUserHandler(userManager *usecases.UserManager, middleware []echo.MiddlewareFunc) *UserHandler {
	return &UserHandler{
		middleware:  middleware,
		userManager: userManager,
	}
}

func (h *UserHandler) RegisterRoutes(ws *echo.Echo) {
	ws.GET("/users/:id", h.getUser, h.middleware...)
	ws.PUT("/users/:id", h.updateUser, h.middleware...)
	ws.DELETE("/users/:id", h.deleteUser, h.middleware...)
}

func (h *UserHandler) getUser(ctx echo.Context) error {
	panic("implement me")
}

func (h *UserHandler) updateUser(ctx echo.Context) error {
	panic("implement me")
}

func (h *UserHandler) deleteUser(ctx echo.Context) error {
	panic("implement me")
}
