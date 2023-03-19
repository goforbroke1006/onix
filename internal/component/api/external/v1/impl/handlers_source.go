package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

func (h handlersImpl) PostSource(ctx echo.Context, params spec.PostSourceParams) error {
	createErr := h.sourceRepo.Create(ctx.Request().Context(), params.Title, domain.SourceType(params.Kind), params.Address)
	if createErr != nil {
		return errors.Wrap(createErr, "can't store source in repository")
	}

	return ctx.String(http.StatusCreated, "OK")
}

func (h handlersImpl) GetSource(ctx echo.Context) error {
	sourcesList, loadErr := h.sourceRepo.GetAll(ctx.Request().Context())
	if loadErr != nil {
		return errors.Wrap(loadErr, "can't get sources list")
	}

	response := make([]spec.Source, 0, len(sourcesList))

	for _, src := range sourcesList {
		response = append(response, spec.Source{
			Id:      src.ID,
			Kind:    spec.SourceKind(src.Kind),
			Address: src.Address,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
