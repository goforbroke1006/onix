package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	apiSpec "github.com/goforbroke1006/onix/api/stub_prometheus"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewServer creates new server's handlers implementations instance.
func NewServer(logger log.Logger) *server { // nolint:revive,golint
	return &server{
		validator: validator{},
		logger:    logger,
	}
}

var _ apiSpec.ServerInterface = &server{} // nolint:exhaustivestruct

type server struct {
	validator validator
	logger    log.Logger
}

func (s server) GetHealthz(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "ok")

	return errors.Wrap(err, "write to echo context failed")
}

func (s server) GetQuery(ctx echo.Context) error {
	err := ctx.String(http.StatusOK, "implement me")

	return errors.Wrap(err, "write to echo context failed")
}

func (s server) GetQueryRange(ctx echo.Context, params apiSpec.GetQueryRangeParams) error {
	if err := s.validator.GetQueryRange(params); err != nil {
		s.logger.WithErr(err).Warn("invalid request")
		err := ctx.NoContent(http.StatusBadRequest)

		return errors.Wrap(err, "write to echo context failed")
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

	queryRangeResult := apiSpec.QueryRangeResult{ // nolint:exhaustivestruct
		Metric: apiSpec.QueryRangeResult_Metric{}, // nolint:exhaustivestruct
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
