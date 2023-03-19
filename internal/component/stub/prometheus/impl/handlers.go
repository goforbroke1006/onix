package impl

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/goforbroke1006/onix/internal/component/stub/prometheus"
	"github.com/goforbroke1006/onix/internal/component/stub/prometheus/spec"
)

// NewHandlers creates new handlers's handlers implementations instance.
func NewHandlers() spec.ServerInterface {
	return &handlers{}
}

var _ spec.ServerInterface = (*handlers)(nil)

type handlers struct{}

func (h handlers) GetQuery(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "implement me")

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetQueryRange(ctx echo.Context, params spec.GetQueryRangeParams) error {
	if errs := params.Validate(); errs != nil {
		zap.L().Error("invalid request", zap.Any("errs", errs))
		return ctx.NoContent(http.StatusBadRequest)
	}

	var (
		start   time.Time
		stop    time.Time
		step    time.Duration
		timeout time.Duration
	)

	start, _ = prometheus.CanParseTime(params.Start)
	stop, _ = prometheus.CanParseTime(params.End)
	step, _ = prometheus.CanParseDuration(params.Step)
	timeout, _ = prometheus.CanParseDuration(string(params.Timeout))

	const oneDay = 24 * time.Hour
	if stop.Sub(start) > oneDay {
		stop = start.Add(oneDay)
	}

	resp := spec.QueryRangeResponse{
		Status: spec.StatusSuccess,
		Data: spec.QueryRangeData{
			ResultType: spec.QueryRangeDataResultTypeMatrix,
			Result:     nil,
		},
	}

	queryRangeResult := spec.QueryRangeResult{ //nolint:exhaustivestruct
		Metric: spec.QueryRangeResult_Metric{}, //nolint:exhaustivestruct
	}

	ctx2, cancel2 := context.WithTimeout(ctx.Request().Context(), timeout)

	go func() {
		defer cancel2()

		gen := prometheus.NewFakeMetricsIdempotentGenerator()
		for _, si := range gen.Load(params.Query, start, stop, step) {
			var (
				timestamp = float64(si.Timestamp)
				val       = fmt.Sprintf("%f", si.Value)
			)

			queryRangeResult.Values = append(queryRangeResult.Values, []interface{}{timestamp, val})
		}

		resp.Data.Result = append(resp.Data.Result, queryRangeResult)
	}()

	<-ctx2.Done()

	return ctx.JSON(http.StatusOK, resp)
}
