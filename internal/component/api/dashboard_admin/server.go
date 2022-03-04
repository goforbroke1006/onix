package dashboard_admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-admin"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

func NewServer(
	serviceRepo domain.ServiceRepository,
	releaseRepo domain.ReleaseRepository,
	sourceRepo domain.SourceRepository,
	criteriaRepo domain.CriteriaRepository,
	logger log.Logger,
) *server {
	return &server{
		serviceRepo:  serviceRepo,
		releaseRepo:  releaseRepo,
		sourceRepo:   sourceRepo,
		criteriaRepo: criteriaRepo,
		logger:       logger,
	}
}

var (
	_ apiSpec.ServerInterface = &server{}
)

type server struct {
	serviceRepo  domain.ServiceRepository
	releaseRepo  domain.ReleaseRepository
	sourceRepo   domain.SourceRepository
	criteriaRepo domain.CriteriaRepository
	logger       log.Logger
}

func (s server) GetService(ctx echo.Context) error {
	const (
		defaultReleasesCount = 5
	)

	services, err := s.serviceRepo.GetAll()
	if err != nil {
		return err
	}

	resp := apiSpec.ServicesListResponse{}
	for _, svc := range services {
		service := apiSpec.Service{
			Title:    svc.Title,
			Releases: []string{},
		}

		releases, err := s.releaseRepo.GetNLasts(svc.Title, defaultReleasesCount)
		if err != nil {
			return err
		}

		for _, r := range releases {
			service.Releases = append(service.Releases, r.Name)
		}

		resp = append(resp, service)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s server) GetSource(ctx echo.Context) error {
	sources, err := s.sourceRepo.GetAll()
	if err != nil {
		return err
	}
	resp := make(apiSpec.SourceListResponse, 0, len(sources))
	for _, s := range sources {
		resp = append(resp, apiSpec.Source{
			Id:      s.ID,
			Title:   s.Title,
			Kind:    apiSpec.SourceKind(s.Kind),
			Address: s.Address,
		})
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s server) PostCriteria(ctx echo.Context) error {
	requestBody := apiSpec.CreateCriteriaRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&requestBody); err != nil {
		return err
	}
	criteriaID, err := s.criteriaRepo.Create(
		requestBody.ServiceName, requestBody.Title, requestBody.Selector,
		domain.DynamicDirType(requestBody.ExpectedDir),
		domain.MustParseGroupingIntervalType(string(requestBody.Interval)))
	if err != nil {
		return err
	}

	s.logger.Info("create new criteria")

	resp := apiSpec.CreateResourceResponse{
		NewId:  fmt.Sprintf("%d", criteriaID),
		Status: apiSpec.CreateResourceResponseStatusOk,
	}
	return ctx.JSON(http.StatusOK, resp)
}
