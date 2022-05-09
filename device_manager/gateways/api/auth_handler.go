package api

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/ovidiuz/device_manager/device_manager/jwt"

	"github.com/sirupsen/logrus"

	"github.com/ovidiuz/device_manager/device_manager/domain"
	"github.com/ovidiuz/device_manager/device_manager/usecases"
)

type AuthHandler struct {
	userManager    *usecases.UserManager
	authMiddleware []echo.MiddlewareFunc
	jwtTokenTTL    time.Duration
	httpsEnabled   bool
}

func NewAuthHandler(userManager *usecases.UserManager, jwtTokenTTL time.Duration, httpsEnabled bool) *AuthHandler {
	authHandler := &AuthHandler{
		userManager:  userManager,
		jwtTokenTTL:  jwtTokenTTL,
		httpsEnabled: httpsEnabled,
	}
	authHandler.authMiddleware = []echo.MiddlewareFunc{authHandler.isAuthenticated}

	return authHandler
}

func (h *AuthHandler) RegisterRoutes(ws *echo.Echo) {
	ws.POST("/auth/register", h.register)
	ws.POST("/auth/login", h.login)
}

func (h *AuthHandler) GetAuthMiddleware() []echo.MiddlewareFunc {
	return h.authMiddleware
}

func (h *AuthHandler) register(ctx echo.Context) (apiErr error) {
	registerRequest := &domain.RegisterRequest{}
	if err := ctx.Bind(registerRequest); err != nil {
		logrus.WithError(err).Debug("invalid or missing request body")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := registerRequest.Validate(); err != nil {
		logrus.WithError(err).Debug("validation failed for request body")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userManager.GetUserByEmail(context.TODO(), registerRequest.Email)
	if err == nil && user != nil {
		logrus.Debugf("user with email=%s already exists", registerRequest.Email)
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	} else if err != nil && err != domain.ErrNotFound {
		logrus.WithError(err).Errorf("could not get user with email=%s", registerRequest.Email)
		return err
	}
	// err == domain.ErrNotFound

	user, err = h.userManager.RegisterUser(context.TODO(), registerRequest)
	if err != nil {
		logrus.WithError(err).Errorf("could not register user with email=%s", registerRequest.Email)
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

func (h *AuthHandler) login(ctx echo.Context) (apiErr error) {
	loginRequest := &domain.LoginRequest{}
	if err := ctx.Bind(loginRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := loginRequest.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := h.userManager.LoginUser(context.TODO(), loginRequest)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     domain.JWT,
		Value:    token,
		Expires:  time.Now().Add(h.jwtTokenTTL),
		MaxAge:   int(h.jwtTokenTTL),
		Secure:   h.httpsEnabled,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}

func (h *AuthHandler) isAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(domain.JWT)
		if err == http.ErrNoCookie {
			return echo.ErrUnauthorized
		} else if err != nil {
			return echo.ErrInternalServerError
		}

		userID, err := jwt.ParseJWT(cookie.Value)
		if err != nil {
			return echo.ErrUnauthorized
		}

		_, err = h.userManager.GetUser(ctx.Request().Context(), userID)
		if err == domain.ErrNotFound {
			return echo.ErrUnauthorized
		}
		if err != nil {
			return echo.ErrInternalServerError
		}

		ctx.Set(domain.UserKey, userID)
		return next(ctx)
	}
}
