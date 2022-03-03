package prometheus

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func NewServer() *server {
	return &server{}
}

var (
	_ ServerInterface = &server{}
)

type server struct {
}

func (s server) GetHealthz(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "ok")
}

func (s server) GetQuery(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s server) GetQueryRange(ctx echo.Context, params GetQueryRangeParams) error {
	// TODO: fix auto-gen validation
	if len(params.Query) == 0 {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var (
		start time.Time
		stop  time.Time
		step  time.Duration
		err   error
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

	resp := QueryRangeResponse{
		Status: StatusSuccess,
		Data: QueryRangeData{
			ResultType: QueryRangeDataResultTypeMatrix,
			Result:     nil,
		},
	}

	queryRangeResult := QueryRangeResult{
		Metric: QueryRangeResult_Metric{},
	}
	gen := fakeMetricsIdempotentGenerator{}
	for _, si := range gen.Load(params.Query, start, stop, step) {
		var (
			timestamp = float64(si.timestamp)
			val       = fmt.Sprintf("%f", si.value)
		)
		queryRangeResult.Values = append(queryRangeResult.Values, []interface{}{timestamp, val})
	}
	resp.Data.Result = append(resp.Data.Result, queryRangeResult)

	return ctx.JSON(http.StatusOK, resp)
}

func (s server) canParseTime(str string) (time.Time, error) {
	var (
		onlyNumbersRegex = regexp.MustCompile(`^[\d]+$`)
		res              time.Time
		err              error
	)
	if onlyNumbersRegex.Match([]byte(str)) {
		unix, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			panic(err)
		}
		res = time.Unix(unix, 0).UTC()
		return res, nil
	}

	res, err = time.Parse(time.RFC3339, str)
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
