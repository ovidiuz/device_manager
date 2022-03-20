package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

func Authorize(obj, act string, enforcer casbin.IEnforcer) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// TODO: add role checking
		return nil
	}
}
