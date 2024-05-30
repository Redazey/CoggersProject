package mw

import (
	"CoggersProject/backend/pkg/service/logger"
	"fmt"
	"runtime"

	"github.com/labstack/echo"
)

func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				funcName, file, line, ok := runtime.Caller(2) // 2 - количество уровней стека, для получения имени функции, вызвавшей панику
				if ok {
					funcDetails := runtime.FuncForPC(funcName)
					logStr := fmt.Sprintf("в ходе выполнения функции %s в файле %s на строке %d возникла паника: %s", funcDetails.Name(), file, line, r)
					logger.Error(logStr)
				} else {
					logStr := fmt.Sprintf("в ходе выполнения функции возника паника: %s", r)
					logger.Error(logStr)
				}
			}
		}()
		return next(c)
	}
}
