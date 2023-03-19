package metricsprovider

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/external/prom"
)

// NewPrometheusMetricsProvider inits new prom data provider.
func NewPrometheusMetricsProvider(address string) *promMetricsProvider { //nolint:revive,golint
	return &promMetricsProvider{
		client: prom.NewClient(address),
	}
}

var _ domain.MetricsProvider = &promMetricsProvider{} //nolint:exhaustivestruct

type promMetricsProvider struct {
	client prom.APIClient
}

const defaultTimeout = 10 * time.Second

func (p promMetricsProvider) LoadSeries(
	ctx context.Context,
	selector string,
	from, till time.Time,
	step time.Duration,
) ([]domain.SeriesItem, error) {
	var (
		resp *prom.QueryRangeResponse
		err  error
	)

	if resp, err = p.client.QueryRange(ctx, selector, from, till, step, defaultTimeout); err != nil {
		return nil, errors.Wrap(err, "wrong response from prom API")
	}

	series := make([]domain.SeriesItem, 0, len(resp.Data.Result[0].Values))

	for _, graphValue := range resp.Data.Result[0].Values {
		f, ok := graphValue[0].(float64)
		if !ok {
			continue
		}

		unix := int64(f)
		moment := time.Unix(unix, 0)

		const bitsSize = 64
		value, _ := strconv.ParseFloat(graphValue[1].(string), bitsSize)

		series = append(series, domain.SeriesItem{
			Timestamp: moment,
			Value:     value,
		})
	}

	return series, nil
}
