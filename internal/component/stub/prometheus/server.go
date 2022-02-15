package prometheus

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
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
	rand.Seed(time.Now().UnixNano())

	resp := QueryRangeResponse{
		Status: "",
		Data: QueryRangeData{
			ResultType: QueryRangeDataResultTypeMatrix,
			Result:     nil,
		},
	}

	var (
		current = s.mustParseTime(params.Start)
		stop    = s.mustParseTime(params.End)
		step    = s.mustParseDuration(params.Step)
	)
	queryRangeResult := QueryRangeResult{
		Metric: QueryRangeResult_Metric{},
	}
	for current.Before(stop) || current.Equal(stop) {
		var (
			timestamp = float64(current.Unix())
			val       = fmt.Sprintf("%f", rand.Float64())
		)
		queryRangeResult.Values = append(queryRangeResult.Values, []interface{}{timestamp, val})

		current = current.Add(step)
	}
	resp.Data.Result = append(resp.Data.Result, queryRangeResult)

	return ctx.JSON(http.StatusOK, resp)
}

func (s server) mustParseTime(str string) time.Time {
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
		res = time.Unix(unix, 0)
		return res
	}

	res, err = time.Parse(time.RFC3339, str)
	if err != nil {
		panic(err)
	}

	return res
}

func (s server) mustParseDuration(str string) time.Duration {
	float, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return time.Duration(float) * time.Second
	}

	duration, err := time.ParseDuration(str)
	if err == nil {
		return duration
	}

	panic("can't parse duration")
}
