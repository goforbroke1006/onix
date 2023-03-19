package impl

import (
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

// NewHandlers creates new handlersImpl instance.
func NewHandlers(
	serviceRepo domain.ServiceRepository,
	sourceRepo domain.SourceRepository,
	criteriaRepo domain.CriteriaRepository,
	releaseRepo domain.ReleaseRepository,
) spec.ServerInterface {
	return &handlersImpl{
		serviceRepo:  serviceRepo,
		sourceRepo:   sourceRepo,
		criteriaRepo: criteriaRepo,
		releaseRepo:  releaseRepo,
	}
}

var _ spec.ServerInterface = (*handlersImpl)(nil)

type handlersImpl struct {
	serviceRepo        domain.ServiceRepository
	sourceRepo         domain.SourceRepository
	criteriaRepo       domain.CriteriaRepository
	releaseRepo        domain.ReleaseRepository
	releaseSvc         domain.ReleaseService
	measurementService domain.MeasurementService
}
