package dashboardadmin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-admin"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewHandlers creates new handlers implementations instance.
func NewHandlers(
	serviceRepo domain.ServiceRepository,
	releaseRepo domain.ReleaseRepository,
	sourceRepo domain.SourceRepository,
	criteriaRepo domain.CriteriaRepository,
	logger log.Logger,
) *handlers { // nolint:revive,golint
	return &handlers{
		serviceRepo:  serviceRepo,
		releaseRepo:  releaseRepo,
		sourceRepo:   sourceRepo,
		criteriaRepo: criteriaRepo,
		logger:       logger,
	}
}

var _ apiSpec.ServerInterface = &handlers{} // nolint:exhaustivestruct

type handlers struct {
	serviceRepo  domain.ServiceRepository
	releaseRepo  domain.ReleaseRepository
	sourceRepo   domain.SourceRepository
	criteriaRepo domain.CriteriaRepository
	logger       log.Logger
}

func (h handlers) GetHealthz(ctx echo.Context) error {
	err := ctx.NoContent(http.StatusOK)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetService(ctx echo.Context) error {
	const (
		defaultReleasesCount = 5
	)

	services, err := h.serviceRepo.GetAll()
	if err != nil {
		return errors.Wrap(err, "can't get services from repository")
	}

	resp := apiSpec.ServicesListResponse{}

	for _, svc := range services {
		service := apiSpec.Service{
			Title:    svc.Title,
			Releases: []string{},
		}

		releases, err := h.releaseRepo.GetNLasts(svc.Title, defaultReleasesCount)
		if err != nil {
			return errors.Wrap(err, "can't get N last releases")
		}

		for _, r := range releases {
			service.Releases = append(service.Releases, r.Name)
		}

		resp = append(resp, service)
	}

	err = ctx.JSON(http.StatusOK, resp)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetSource(ctx echo.Context) error {
	var (
		sources []domain.Source
		err     error
	)

	if sources, err = h.sourceRepo.GetAll(); err != nil {
		return errors.Wrap(err, "can't get sources")
	}

	resp := make(apiSpec.SourceListResponse, 0, len(sources))
	for _, s := range sources {
		resp = append(resp, apiSpec.Source{
			Id:      s.ID,
			Title:   s.Title,
			Kind:    apiSpec.SourceKind(s.Type),
			Address: s.Address,
		})
	}

	err = ctx.JSON(http.StatusOK, resp)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) PostCriteria(ctx echo.Context) error {
	var requestBody apiSpec.CreateCriteriaRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&requestBody); err != nil {
		return errors.Wrap(err, "incorrect post body")
	}

	criteriaID, err := h.criteriaRepo.Create(
		requestBody.ServiceName, requestBody.Title, requestBody.Selector,
		domain.DynamicDirType(requestBody.ExpectedDir),
		domain.MustParseGroupingIntervalType(string(requestBody.Interval)))
	if err != nil {
		return errors.Wrap(err, "can't store criteria in db")
	}

	h.logger.Info("create new criteria")

	resp := apiSpec.CreateResourceResponse{
		NewId:  fmt.Sprintf("%d", criteriaID),
		Status: apiSpec.CreateResourceResponseStatusOk,
	}
	err = ctx.JSON(http.StatusOK, resp)

	return errors.Wrap(err, "write to echo context failed")
}
