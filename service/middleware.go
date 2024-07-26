package service

import (
	"github.com/labstack/echo/v4"

	"github.com/sika365/admin-tools/context"
)

func (svc *Service) UseContext(e *echo.Echo) error {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.Context{Context: c}
			return next(cc)
		}
	})
	return nil
}
