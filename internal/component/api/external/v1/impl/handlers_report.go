package impl

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

func (h handlersImpl) PostReportCompare(ctx echo.Context) error { // nolint:funlen,cyclop
	var params spec.ReportCompareRequest
	if bindErr := ctx.Bind(params); bindErr != nil {
		return ctx.JSON(http.StatusBadRequest, spec.NewErrorResponse(bindErr))
	}

	reqCtx := ctx.Request().Context()

	period, _ := time.ParseDuration(string(params.TimeRange))

	source, sourceErr := h.sourceRepo.Get(reqCtx, params.Source)
	if sourceErr != nil {
		return errors.Wrap(sourceErr, "can't get source")
	}

	releaseOne, err := h.releaseSvc.GetByName(params.Service, params.TagOne)
	if err != nil {
		return errors.Wrap(err, "can't get release one by name")
	}

	releaseTwo, err := h.releaseSvc.GetByName(params.Service, params.TagTwo)
	if err != nil {
		return errors.Wrap(err, "can't get release two by name")
	}

	criteriaList, err := h.criteriaRepo.GetAll(params.Service)
	if err != nil {
		return errors.Wrap(err, "can't get criteria list")
	}

	response := spec.ReportCompareResponse{}

	const defaultSeriesStep = 5 * time.Minute

	for _, criteria := range criteriaList {
		var (
			series1 []domain.MeasurementRow
			series2 []domain.MeasurementRow
			err     error
		)

		if series1, err = h.measurementService.GetOrPull(
			reqCtx,
			source, criteria, releaseOne.StartAt, releaseOne.StartAt.Add(period), defaultSeriesStep,
		); err != nil {
			return errors.Wrap(err, "can't get series for release one")
		}

		if series2, err = h.measurementService.GetOrPull(
			reqCtx,
			source, criteria, releaseTwo.StartAt, releaseTwo.StartAt.Add(period), defaultSeriesStep,
		); err != nil {
			return errors.Wrap(err, "can't get series for release two")
		}

		minLen := len(series1)
		if len(series2) < minLen {
			minLen = len(series2)
		}

		var direction spec.CriteriaReportDirection

		switch criteria.Direction {
		case domain.DynamicDirTypeIncrease:
			direction = spec.CriteriaReportDirectionIncrease
		case domain.DynamicDirTypeDecrease:
			direction = spec.CriteriaReportDirectionDecrease
		case domain.DynamicDirTypeEqual:
			direction = spec.CriteriaReportDirectionEqual
		}

		criteriaReport := spec.CriteriaReport{
			Title:     criteria.Title,
			Selector:  criteria.Selector,
			Graph:     make([]spec.GraphItem, 0, minLen),
			Direction: direction,
		}

		for seriesItemIndex := 0; seriesItemIndex < minLen; seriesItemIndex++ {
			criteriaReport.Graph = append(criteriaReport.Graph, spec.GraphItem{
				T1: series1[seriesItemIndex].Moment.UTC().Unix(),
				V1: series1[seriesItemIndex].Value,

				T2: series2[seriesItemIndex].Moment.UTC().Unix(),
				V2: series2[seriesItemIndex].Value,
			})
		}

		response = append(response, criteriaReport)
	}

	return ctx.JSON(http.StatusOK, response)
}
