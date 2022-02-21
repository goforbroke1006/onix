package dashboard_admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

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
	_ ServerInterface = &server{}
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

	resp := ServicesListResponse{}
	for _, svc := range services {
		service := Service{
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
	resp := make(SourceListResponse, 0, len(sources))
	for _, s := range sources {
		resp = append(resp, Source{
			Id:      s.ID,
			Title:   s.Title,
			Kind:    SourceKind(s.Kind),
			Address: s.Address,
		})
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s server) PostCriteria(ctx echo.Context) error {
	requestBody := CreateCriteriaRequest{}
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

	resp := CreateResourceResponse{
		NewId:  fmt.Sprintf("%d", criteriaID),
		Status: CreateResourceResponseStatusOk,
	}
	return ctx.JSON(http.StatusOK, resp)
}
