package api

import (
	"context"
	"net/http"

	"github.com/casbin/casbin/v2"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/ovidiuz/device_manager/domain"
	"github.com/ovidiuz/device_manager/usecases"
)

type AuthHandler struct {
	userManager *usecases.UserManager
	enforcer    casbin.IEnforcer
}

func NewAuthHandler(userManager *usecases.UserManager, enforcer casbin.IEnforcer) *AuthHandler {
	return &AuthHandler{
		userManager: userManager,
		enforcer:    enforcer,
	}
}

func (h *AuthHandler) RegisterRoutes(ws *echo.Echo) {
	ws.POST("/auth/register", h.register)
	ws.POST("/auth/login", h.login)
}

func (h *AuthHandler) register(ctx echo.Context) (apiErr error) {
	var handlerErr *domain.HandlerErr
	defer func() {
		if handlerErr != nil {
			apiErr = ctx.JSON(handlerErr.Code, handlerErr)
		}
	}()

	registerRequest := &domain.RegisterRequest{}
	if err := ctx.Bind(registerRequest); err != nil {
		logrus.WithError(err).Debug("invalid or missing request body")
		handlerErr = domain.NewHandlerErr(http.StatusBadRequest, domain.InvalidRequestBody)
		return
	}
	if err := registerRequest.Validate(); err != nil {
		logrus.WithError(err).Debug("validation failed for request body")
		handlerErr = domain.NewHandlerErr(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userManager.RegisterUser(context.TODO(), registerRequest)
	if err != nil {
		logrus.WithError(err).Errorf("could not register user with email=%s", user.Email)
		handlerErr = domain.NewHandlerErr(http.StatusInternalServerError, domain.InternalServerError)
		return
	}

	return ctx.JSON(http.StatusOK, user)
}

func (h *AuthHandler) login(ctx echo.Context) (apiErr error) {
	var handlerErr *domain.HandlerErr
	defer func() {
		if handlerErr != nil {
			apiErr = ctx.JSON(handlerErr.Code, handlerErr)
		}
	}()

	loginRequest := &domain.LoginRequest{}
	if err := ctx.Bind(loginRequest); err != nil {
		handlerErr = domain.NewHandlerErr(http.StatusBadRequest, "invalid request body")
		return
	}
	if err := loginRequest.Validate(); err != nil {
		handlerErr = domain.NewHandlerErr(http.StatusBadRequest, "invalid request body")
		return
	}
	panic("implement me")
}
