package system

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	apiSpec "github.com/goforbroke1006/onix/api/system"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewHandlers creates new handlers's handlers implementations instance.
func NewHandlers(
	serviceRepo domain.ServiceRepository,
	sourceRepo domain.SourceRepository,
	criteriaRepo domain.CriteriaRepository,
	releaseRepo domain.ReleaseRepository,
	logger log.Logger,
) *handlers { // nolint:revive,golint
	return &handlers{
		serviceRepo:  serviceRepo,
		sourceRepo:   sourceRepo,
		criteriaRepo: criteriaRepo,
		releaseRepo:  releaseRepo,
		logger:       logger,
	}
}

var _ apiSpec.ServerInterface = &handlers{} // nolint:exhaustivestruct

type handlers struct {
	serviceRepo  domain.ServiceRepository
	sourceRepo   domain.SourceRepository
	criteriaRepo domain.CriteriaRepository
	releaseRepo  domain.ReleaseRepository
	logger       log.Logger
}

func (h handlers) GetHealthz(ctx echo.Context) error {
	err := ctx.NoContent(http.StatusOK)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetRegisterSource(ctx echo.Context, params apiSpec.GetRegisterSourceParams) error {
	if _, err := h.sourceRepo.Create(params.Title, domain.SourceType(params.Kind), params.Address); err != nil {
		return errors.Wrap(err, "can't store source in repository")
	}

	response := apiSpec.RegisterSourceResponse{
		Status: apiSpec.RegisterSourceResponseStatusOk,
	}
	err := ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetRegisterCriteria(ctx echo.Context, params apiSpec.GetRegisterCriteriaParams) error {
	if _, err := h.criteriaRepo.Create(
		params.Service,
		params.Title,
		params.Selector,
		domain.DynamicDirType(params.ExpectedDir),
		domain.MustParseGroupingIntervalType(params.GroupingInterval),
	); err != nil {
		return errors.Wrap(err, "can't store criteria in repository")
	}

	response := apiSpec.RegisterCriteriaResponse{
		Status: apiSpec.RegisterCriteriaResponseStatusOk,
	}
	err := ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetRegisterRelease(ctx echo.Context, params apiSpec.GetRegisterReleaseParams) error {
	startAt := time.Now().UTC()
	if params.StartAt != nil {
		startAt = time.Unix(*params.StartAt, 0).UTC()
	}

	if err := h.serviceRepo.Store(params.ServiceName); err != nil {
		return errors.Wrap(err, "can't store service in repository")
	}

	if err := h.releaseRepo.Store(params.ServiceName, params.ReleaseName, startAt); err != nil {
		return errors.Wrap(err, "can't get release")
	}

	response := apiSpec.RegisterReleaseResponse{
		Status: apiSpec.RegisterReleaseResponseStatusOk,
	}
	err := ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}
