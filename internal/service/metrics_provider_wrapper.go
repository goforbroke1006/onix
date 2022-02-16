package service

import (
	"fmt"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/service/metrics_provider"
)

func NewMetricsProvider(source domain.Source) domain.MetricsProvider {
	switch source.Kind {
	case domain.SourceTypePrometheus:
		return metrics_provider.NewPrometheusMetricsProvider(source.Address)
	case domain.SourceTypeInfluxDB:
		return metrics_provider.NewInfluxDBMetricsProvider()
	default:
		panic(fmt.Errorf("unexpected metrics provider type: %s", source.Kind))
	}
}
