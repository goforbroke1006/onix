package metrics_provider

import (
	"time"

	"github.com/goforbroke1006/onix/domain"
)

func NewInfluxDBMetricsProvider() *influxDbMetricsProvider {
	return &influxDbMetricsProvider{}
}

var (
	_ domain.MetricsProvider = &influxDbMetricsProvider{}
)

type influxDbMetricsProvider struct {
}

func (p influxDbMetricsProvider) LoadSeries(
	selector string,
	from, till time.Time,
	step time.Duration,
) ([]domain.SeriesItem, error) {
	//TODO implement me
	panic("implement me")
}
