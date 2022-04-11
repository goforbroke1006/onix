package dashboardmain

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-main"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewHandlers creates new handlers's handlers implementations instance.
func NewHandlers(
	serviceRepo domain.ServiceRepository,
	releaseSvc domain.ReleaseService,
	sourceRepo domain.SourceRepository,
	criteriaRepo domain.CriteriaRepository,
	measurementService domain.MeasurementService,
	logger log.Logger,
) *handlers { // nolint:revive,golint
	return &handlers{
		serviceRepo:        serviceRepo,
		releaseSvc:         releaseSvc,
		sourceRepo:         sourceRepo,
		criteriaRepo:       criteriaRepo,
		measurementService: measurementService,
		logger:             logger,
	}
}

var _ apiSpec.ServerInterface = &handlers{} // nolint:exhaustivestruct

type handlers struct {
	serviceRepo        domain.ServiceRepository
	releaseSvc         domain.ReleaseService
	sourceRepo         domain.SourceRepository
	criteriaRepo       domain.CriteriaRepository
	measurementService domain.MeasurementService
	logger             log.Logger
}

func (h handlers) GetHealthz(ctx echo.Context) error {
	err := ctx.NoContent(http.StatusOK)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetService(ctx echo.Context) error {
	services, err := h.serviceRepo.GetAll()
	if err != nil {
		return errors.Wrap(err, "can't get services list")
	}

	response := make([]apiSpec.Service, 0, len(services))
	for _, svc := range services {
		response = append(response, apiSpec.Service{Title: svc.Title})
	}

	err = ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetSource(ctx echo.Context) error {
	sourcesList, err := h.sourceRepo.GetAll()
	if err != nil {
		return errors.Wrap(err, "can't get sources list")
	}

	response := make([]apiSpec.Source, 0, len(sourcesList))

	for _, src := range sourcesList {
		response = append(response, apiSpec.Source{
			Id:      src.ID,
			Title:   src.Title,
			Kind:    apiSpec.SourceKind(src.Type),
			Address: src.Address,
		})
	}

	err = ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetRelease(ctx echo.Context, params apiSpec.GetReleaseParams) error {
	ranges, err := h.releaseSvc.GetAll(params.Service)
	if err != nil {
		return errors.Wrap(err, "can't get releases list")
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

	err = ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetCompare(ctx echo.Context, params apiSpec.GetCompareParams) error { // nolint:funlen,cyclop
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

	releaseOne, err := h.releaseSvc.GetByName(params.Service, params.ReleaseOneTitle)
	if err != nil {
		return errors.Wrap(err, "can't get release one by name")
	}

	if releaseOneStart.Before(releaseOne.StartAt) {
		h.logger.Errorf("%d before %d", params.ReleaseOneStart, releaseOne.StartAt.Unix())

		return errors.Wrap(ErrWrongTimeRange, "release 1 wrong time")
	}

	releaseTwo, err := h.releaseSvc.GetByName(params.Service, params.ReleaseTwoTitle)
	if err != nil {
		return errors.Wrap(err, "can't get release two by name")
	}

	if releaseTwoStart.Before(releaseTwo.StartAt) {
		h.logger.Errorf("%d before %d", params.ReleaseTwoStart, releaseTwo.StartAt.Unix())

		return errors.Wrap(ErrWrongTimeRange, "release 2 wrong time")
	}

	criteriaList, err := h.criteriaRepo.GetAll(params.Service)
	if err != nil {
		return errors.Wrap(err, "can't get criteria list")
	}

	response := apiSpec.CompareResponse{ // nolint:exhaustivestruct
		Service:    params.Service,
		ReleaseOne: params.ReleaseOneTitle,
		ReleaseTwo: params.ReleaseTwoTitle,
	}

	sourceOne, err := h.sourceRepo.Get(params.ReleaseOneSourceId)
	if err != nil {
		return errors.Wrap(err, "can't get source #1")
	}

	sourceTwo, err := h.sourceRepo.Get(params.ReleaseTwoSourceId)
	if err != nil {
		return errors.Wrap(err, "can't get source #2")
	}

	const defaultSeriesStep = 5 * time.Minute

	for _, criteria := range criteriaList {
		var (
			series1 []domain.MeasurementRow
			series2 []domain.MeasurementRow
			err     error
		)

		if series1, err = h.measurementService.GetOrPull(context.TODO(),
			*sourceOne, criteria, releaseOneStart, releaseOneStop, defaultSeriesStep,
		); err != nil {
			return errors.Wrap(err, "can't get series for release one")
		}

		if series2, err = h.measurementService.GetOrPull(context.TODO(),
			*sourceTwo, criteria, releaseTwoStart, releaseTwoStop, defaultSeriesStep,
		); err != nil {
			return errors.Wrap(err, "can't get series for release two")
		}

		minLen := len(series1)
		if len(series2) < minLen {
			minLen = len(series2)
		}

		criteriaReport := apiSpec.CriteriaReport{
			Title:    criteria.Title,
			Selector: criteria.Selector,
			Graph:    make([]apiSpec.GraphItem, 0, minLen),
		}

		for seriesItemIndex := 0; seriesItemIndex < minLen; seriesItemIndex++ {
			criteriaReport.Graph = append(criteriaReport.Graph, apiSpec.GraphItem{
				T1: series1[seriesItemIndex].Moment.UTC().Format(layout),
				V1: series1[seriesItemIndex].Value,

				T2: series2[seriesItemIndex].Moment.UTC().Format(layout),
				V2: series2[seriesItemIndex].Value,
			})
		}

		response.Reports = append(response.Reports, criteriaReport)
	}

	err = ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}
