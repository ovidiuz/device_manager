package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ovidiuz/device_manager/domain"
	"github.com/ovidiuz/device_manager/usecases"
)

type AuthHandler struct {
	userManager *usecases.UserManager
}

func NewAuthHandler(userManager *usecases.UserManager) *AuthHandler {
	return &AuthHandler{
		userManager: userManager,
	}
}

func (h *AuthHandler) RegisterRoutes(ws *echo.Echo) {
	ws.POST("/auth/register", h.register)
	ws.POST("/auth/login", h.login)
}

func (h *AuthHandler) register(ctx echo.Context) error {
	registerRequest := &domain.RegisterRequest{}
	if err := ctx.Bind(registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := registerRequest.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	panic("implement me")
}

func (h *AuthHandler) login(ctx echo.Context) error {
	panic("implement me")
}
