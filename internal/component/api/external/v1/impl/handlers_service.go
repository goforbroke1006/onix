package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	apiSpec "github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

func (h handlersImpl) GetService(ctx echo.Context) error {
	services, err := h.serviceRepo.GetAll(ctx.Request().Context())
	if err != nil {
		return errors.Wrap(err, "can't get services list")
	}

	response := make([]apiSpec.Service, 0, len(services))
	for _, svc := range services {
		response = append(response, apiSpec.Service{Title: svc.ID})
	}

	return ctx.JSON(http.StatusOK, response)
}
