package metricsprovider

import (
	"context"
	"time"

	"github.com/goforbroke1006/onix/domain"
)

// NewInfluxDBMetricsProvider inits new influx data provider.
func NewInfluxDBMetricsProvider() *influxDBMetricsProvider { // nolint:revive,golint
	return &influxDBMetricsProvider{}
}

var _ domain.MetricsProvider = &influxDBMetricsProvider{}

type influxDBMetricsProvider struct{}

func (p influxDBMetricsProvider) LoadSeries(
	ctx context.Context,
	selector string,
	from, till time.Time,
	step time.Duration,
) ([]domain.SeriesItem, error) {
	panic("implement me")
}
