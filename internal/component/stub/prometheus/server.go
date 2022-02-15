package prometheus

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
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

func (s server) GetQuery(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s server) GetQueryRange(ctx echo.Context, params GetQueryRangeParams) error {
	rand.Seed(time.Now().UnixNano())

	resp := QueryRangeResponse{}
	resp.ResultType = QueryRangeResponseResultTypeMatrix

	var (
		current = time.Unix(params.Start.(int64), 0)
		stop    = time.Unix(params.End.(int64), 0)
		step, _ = time.ParseDuration(params.Step.(string))
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
	resp.Result = append(resp.Result, queryRangeResult)

	return ctx.JSON(http.StatusOK, resp)
}
