package metricsprovider

import (
	"strconv"
	"time"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/external/prom"
)

// NewPrometheusMetricsProvider inits new prom data provider
func NewPrometheusMetricsProvider(address string) *promMetricsProvider { // nolint:revive,golint
	return &promMetricsProvider{
		client: prom.NewClient(address),
	}
}

var (
	_ domain.MetricsProvider = &promMetricsProvider{}
)

type promMetricsProvider struct {
	client prom.APIClient
}

func (p promMetricsProvider) LoadSeries(
	selector string,
	from, till time.Time,
	step time.Duration,
) ([]domain.SeriesItem, error) {
	resp, err := p.client.QueryRange(selector, from, till, step, 10*time.Second)
	if err != nil {
		return nil, err
	}
	series := make([]domain.SeriesItem, 0, len(resp.Data.Result[0].Values))
	for _, gv := range resp.Data.Result[0].Values {
		f, ok := gv[0].(float64)
		if !ok {
			continue
		}
		unix := int64(f)
		moment := time.Unix(unix, 0)
		value, _ := strconv.ParseFloat(gv[1].(string), 64)

		series = append(series, domain.SeriesItem{
			Timestamp: moment,
			Value:     value,
		})
	}
	return series, nil
}
