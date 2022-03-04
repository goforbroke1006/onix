package dashboard_main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-main"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

func NewServer(
	serviceRepo domain.ServiceRepository,
	releaseSvc domain.ReleaseService,
	sourceRepo domain.SourceRepository,
	criteriaRepo domain.CriteriaRepository,
	measurementRepo domain.MeasurementRepository,
	logger log.Logger,
) *server {
	return &server{
		serviceRepo:     serviceRepo,
		releaseSvc:      releaseSvc,
		sourceRepo:      sourceRepo,
		criteriaRepo:    criteriaRepo,
		measurementRepo: measurementRepo,
		logger:          logger,
	}
}

var (
	_ apiSpec.ServerInterface = &server{}
)

type server struct {
	serviceRepo     domain.ServiceRepository
	releaseSvc      domain.ReleaseService
	sourceRepo      domain.SourceRepository
	criteriaRepo    domain.CriteriaRepository
	measurementRepo domain.MeasurementRepository
	logger          log.Logger
}

func (s server) GetHealthz(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (s server) GetService(ctx echo.Context) error {
	services, err := s.serviceRepo.GetAll()
	if err != nil {
		return err
	}

	response := make([]apiSpec.Service, 0, len(services))
	for _, svc := range services {
		response = append(response, apiSpec.Service{Title: svc.Title})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s server) GetSource(ctx echo.Context) error {
	sourcesList, err := s.sourceRepo.GetAll()
	if err != nil {
		return err
	}

	response := make([]apiSpec.Source, 0, len(sourcesList))

	for _, src := range sourcesList {
		response = append(response, apiSpec.Source{
			Id:      src.ID,
			Title:   src.Title,
			Kind:    apiSpec.SourceKind(src.Kind),
			Address: src.Address,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s server) GetRelease(ctx echo.Context, params apiSpec.GetReleaseParams) error {
	ranges, err := s.releaseSvc.GetAll(params.Service)
	if err != nil {
		return err
	}

	response := make([]apiSpec.Release, 0, len(ranges))
	for _, r := range ranges {
		response = append(response, apiSpec.Release{
			Id:    r.ID,
			Title: r.Name,
			From:  r.StartAt.Unix(),
			Till:  r.StopAt.Unix(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s server) GetCompare(ctx echo.Context, params apiSpec.GetCompareParams) error {
	const layout = "2006-01-02 15:04:05"

	var (
		releaseOneStart = time.Unix(params.ReleaseOneStart, 0)
		releaseTwoStart = time.Unix(params.ReleaseTwoStart, 0)
		period, _       = time.ParseDuration(string(params.Period))
	)
	var (
		releaseOneStop = releaseOneStart.Add(period)
		releaseTwoStop = releaseTwoStart.Add(period)
	)

	releaseOne, err := s.releaseSvc.GetByName(params.Service, params.ReleaseOneTitle)
	if err != nil {
		return err
	}
	if releaseOneStart.Before(releaseOne.StartAt) {
		return fmt.Errorf("%d before %d", params.ReleaseOneStart, releaseOne.StartAt.Unix())
	}

	releaseTwo, err := s.releaseSvc.GetByName(params.Service, params.ReleaseTwoTitle)
	if err != nil {
		return err
	}
	if releaseTwoStart.Before(releaseTwo.StartAt) {
		return err
	}

	criteriaList, err := s.criteriaRepo.GetAll(params.Service)
	if err != nil {
		return err
	}

	response := apiSpec.CompareResponse{
		Service:    params.Service,
		ReleaseOne: params.ReleaseOneTitle,
		ReleaseTwo: params.ReleaseTwoTitle,
	}

	for _, cr := range criteriaList {

		m1, err := s.measurementRepo.GetBy(params.ReleaseOneSourceId, cr.ID, releaseOneStart, releaseOneStop)
		if err != nil {
			return err
		}
		m2, err := s.measurementRepo.GetBy(params.ReleaseTwoSourceId, cr.ID, releaseTwoStart, releaseTwoStop)
		if err != nil {
			return err
		}

		minLen := len(m1)
		if len(m2) < minLen {
			minLen = len(m2)
		}

		criteriaReport := apiSpec.CriteriaReport{
			Title:    cr.Title,
			Selector: cr.Selector,
			Graph:    make([]apiSpec.GraphItem, 0, minLen),
		}

		for vi := 0; vi < minLen; vi++ {
			criteriaReport.Graph = append(criteriaReport.Graph, apiSpec.GraphItem{
				T1: m1[vi].Moment.UTC().Format(layout),
				V1: m1[vi].Value,

				T2: m2[vi].Moment.UTC().Format(layout),
				V2: m2[vi].Value,
			})
		}

		response.Reports = append(response.Reports, criteriaReport)
	}

	return ctx.JSON(http.StatusOK, response)
}
