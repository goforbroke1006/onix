package prometheus

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	apiSpec "github.com/goforbroke1006/onix/api/stub_prometheus"
)

// NewServer creates new server's handlers implementations instance
func NewServer() *server { // nolint:revive,golint
	return &server{}
}

var (
	_ apiSpec.ServerInterface = &server{}
)

type server struct {
}

func (s server) GetHealthz(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "ok")
}

func (s server) GetQuery(_ echo.Context) error {
	// TODO implement me
	panic("implement me")
}

func (s server) GetQueryRange(ctx echo.Context, params apiSpec.GetQueryRangeParams) error {
	// TODO: fix auto-gen validation
	if len(params.Query) == 0 {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var (
		start   time.Time
		stop    time.Time
		step    time.Duration
		timeout time.Duration
		err     error
	)

	// TODO: fix auto-gen validation
	if start, err = s.canParseTime(params.Start); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	// TODO: fix auto-gen validation
	if stop, err = s.canParseTime(params.End); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	// TODO: fix auto-gen validation
	if step, _ = s.canParseDuration(params.Step); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	// TODO: fix auto-gen validation
	if timeout, _ = s.canParseDuration(string(params.Timeout)); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if stop.Before(start) {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if stop.Sub(start) > 24*time.Hour {
		stop = start.Add(24 * time.Hour)
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
	return ctx.JSON(http.StatusOK, resp)
}

func (s server) canParseTime(str string) (time.Time, error) {
	var (
		onlyNumbersRegex = regexp.MustCompile(`^[\d]+$`)
	)
	if onlyNumbersRegex.Match([]byte(str)) {
		unix, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			panic(err)
		}
		res := time.Unix(unix, 0).UTC()
		return res, nil
	}

	res, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, err
	}

	return res, nil
}

func (s server) canParseDuration(str string) (time.Duration, error) {
	float, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return time.Duration(int(float*math.Pow(10, 9))) * time.Nanosecond, nil
	}

	duration, err := time.ParseDuration(str)
	if err == nil {
		return duration, nil
	}

	return 0, fmt.Errorf("can't parse duration: %s", str)
}
