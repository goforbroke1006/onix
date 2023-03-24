package impl

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

func (h handlersImpl) PostCriteria(ctx echo.Context, params spec.PostCriteriaParams) error {
	interval, intParseErr := time.ParseDuration(params.Interval)
	if intParseErr != nil {
		return intParseErr
	}

	if _, err := h.criteriaRepo.Create(
		ctx.Request().Context(),
		params.Service,
		params.Title,
		params.Selector,
		domain.DirectionType(params.Direction),
		interval,
	); err != nil {
		return errors.Wrap(err, "can't store criteria in repository")
	}

	return ctx.NoContent(http.StatusOK)
}
