package metricsprovider

import (
	"context"
	"time"

	"github.com/goforbroke1006/onix/domain"
)

// NewInfluxDBMetricsProvider inits new influx data provider.
func NewInfluxDBMetricsProvider() domain.MetricsProvider {
	return &influxDBMetricsProvider{}
}

var _ domain.MetricsProvider = (*influxDBMetricsProvider)(nil)

type influxDBMetricsProvider struct{}

func (p influxDBMetricsProvider) LoadSeries(
	ctx context.Context,
	selector string,
	from, till time.Time,
	step time.Duration,
) ([]domain.SeriesItem, error) {
	_ = ctx
	_ = selector
	_, _ = from, till
	_ = step

	panic("implement me")
}
