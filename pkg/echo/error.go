package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/onix/pkg/log"
)

// ErrorHandler is middleware to create echo error handler.
func ErrorHandler(logger log.Logger) func(err error, ctx echo.Context) {
	return func(err error, ctx echo.Context) {
		type errResponse struct {
			Error string
		}

		logger.WithErr(err).Info("api handler error")

		_ = ctx.JSON(http.StatusInternalServerError, errResponse{Error: err.Error()})
	}
}
