package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/ovidiuz/device_manager/usecases"
)

type UserHandler struct {
	userManager *usecases.UserManager
	enforcer    casbin.IEnforcer
}

func NewUserHandler(userManager *usecases.UserManager, enforcer casbin.IEnforcer) *UserHandler {
	return &UserHandler{
		userManager: userManager,
		enforcer:    enforcer,
	}
}

func (u *UserHandler) RegisterRoutes(ws *echo.Echo) {
	ws.GET("/users/:id", u.getUser)
	ws.PUT("/users/:id", u.updateUser)
	ws.DELETE("/users/:id", u.deleteUser)
}

func (u *UserHandler) getUser(ctx echo.Context) error {
	panic("implement me")
}

func (u *UserHandler) updateUser(ctx echo.Context) error {
	panic("implement me")
}

func (u *UserHandler) deleteUser(ctx echo.Context) error {
	panic("implement me")
}
