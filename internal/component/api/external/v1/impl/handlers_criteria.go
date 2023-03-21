package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

func (h handlersImpl) PostCriteria(ctx echo.Context, params spec.PostCriteriaParams) error {
	if _, err := h.criteriaRepo.Create(
		ctx.Request().Context(),
		params.Service,
		params.Title,
		params.Selector,
		domain.DirectionType(params.Direction),
		domain.MustParseGroupingIntervalType(params.Interval),
	); err != nil {
		return errors.Wrap(err, "can't store criteria in repository")
	}

	return ctx.NoContent(http.StatusOK)
}
