package mw

import (
	"CoggersProject/backend/pkg/service/logger"
	"fmt"

	"github.com/labstack/echo"
)

func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				logStr := fmt.Sprintf("в ходе выполнения функции возника паника: %s", r)
				logger.Error(logStr)
			}
		}()
		return next(c)
	}
}
