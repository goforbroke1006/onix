package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ErrorHandler is middleware to create echo error handler.
func ErrorHandler() func(err error, ctx echo.Context) {
	return func(err error, ctx echo.Context) {
		type errResponse struct {
			Error string
		}

		zap.L().Error("handle request fail", zap.Error(err))

		_ = ctx.JSON(http.StatusInternalServerError, errResponse{Error: err.Error()})
	}
}
