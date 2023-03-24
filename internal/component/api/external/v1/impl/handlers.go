package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

// NewHandlers creates new handlersImpl instance.
func NewHandlers(
	serviceRepo domain.ServiceStorage,
	sourceRepo domain.SourceStorage,
	criteriaRepo domain.CriteriaStorage,
	releaseRepo domain.ReleaseStorage,
	releaseSvc domain.ReleaseService,
) spec.ServerInterface {
	return &handlersImpl{
		serviceRepo:  serviceRepo,
		sourceRepo:   sourceRepo,
		criteriaRepo: criteriaRepo,
		releaseRepo:  releaseRepo,
		releaseSvc:   releaseSvc,
	}
}

var _ spec.ServerInterface = (*handlersImpl)(nil)

type handlersImpl struct {
	serviceRepo        domain.ServiceStorage
	sourceRepo         domain.SourceStorage
	criteriaRepo       domain.CriteriaStorage
	releaseRepo        domain.ReleaseStorage
	releaseSvc         domain.ReleaseService
	measurementService domain.MeasurementService
}

func (h handlersImpl) GetPing(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}
