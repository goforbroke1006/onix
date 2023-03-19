package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	apiSpec "github.com/goforbroke1006/onix/api/stub_prometheus"
)

// NewHandlers creates new handlers's handlers implementations instance.
func NewHandlers() *handlers { //nolint:revive,golint
	return &handlers{
		validator: validator{},
	}
}

var _ apiSpec.ServerInterface = &handlers{} //nolint:exhaustivestruct

type handlers struct {
	validator validator
}

func (h handlers) GetHealthz(ctx echo.Context) error {
	err := ctx.NoContent(http.StatusOK)

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetQuery(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "implement me")

	return errors.Wrap(err, "write to echo context failed")
}

func (h handlers) GetQueryRange(ctx echo.Context, params apiSpec.GetQueryRangeParams) error {
	if err := h.validator.GetQueryRange(params); err != nil {
		zap.L().Error("invalid request", zap.Error(err))
		return ctx.NoContent(http.StatusBadRequest)
	}

	var (
		start   time.Time
		stop    time.Time
		step    time.Duration
		timeout time.Duration
		err     error
	)

	start, _ = canParseTime(params.Start)
	stop, _ = canParseTime(params.End)
	step, _ = canParseDuration(params.Step)
	timeout, _ = canParseDuration(string(params.Timeout))

	const oneDay = 24 * time.Hour
	if stop.Sub(start) > oneDay {
		stop = start.Add(oneDay)
	}

	resp := apiSpec.QueryRangeResponse{
		Status: apiSpec.StatusSuccess,
		Data: apiSpec.QueryRangeData{
			ResultType: apiSpec.QueryRangeDataResultTypeMatrix,
			Result:     nil,
		},
	}

	queryRangeResult := apiSpec.QueryRangeResult{ //nolint:exhaustivestruct
		Metric: apiSpec.QueryRangeResult_Metric{}, //nolint:exhaustivestruct
	}

	ctx2, cancel2 := context.WithTimeout(ctx.Request().Context(), timeout)

	go func() {
		defer cancel2()

		gen := fakeMetricsIdempotentGenerator{}
		for _, si := range gen.Load(params.Query, start, stop, step) {
			var (
				timestamp = float64(si.timestamp)
				val       = fmt.Sprintf("%f", si.value)
			)

			queryRangeResult.Values = append(queryRangeResult.Values, []interface{}{timestamp, val})
		}

		resp.Data.Result = append(resp.Data.Result, queryRangeResult)
	}()

	<-ctx2.Done()

	err = ctx.JSON(http.StatusOK, resp)

	return errors.Wrap(err, "write to echo context failed")
}
