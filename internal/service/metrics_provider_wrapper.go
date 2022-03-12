package service

import (
	"fmt"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/service/metricsprovider"
)

// NewMetricsProvider inits metrics provider from domain.Source
func NewMetricsProvider(source domain.Source) domain.MetricsProvider {
	switch source.Kind {
	case domain.SourceTypePrometheus:
		return metricsprovider.NewPrometheusMetricsProvider(source.Address)
	case domain.SourceTypeInfluxDB:
		return metricsprovider.NewInfluxDBMetricsProvider()
	default:
		panic(fmt.Errorf("unexpected metrics provider type: %s", source.Kind))
	}
}
