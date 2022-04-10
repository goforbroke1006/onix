package service

import (
	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/service/metricsprovider"
)

// ErrUnexpectedProviderType is specific error.
var ErrUnexpectedProviderType = errors.New("unexpected metrics provider type")

// NewMetricsProvider inits metrics provider from domain.Source instance.
func NewMetricsProvider(source domain.Source) domain.MetricsProvider { // nolint:ireturn
	switch source.Type {
	case domain.SourceTypePrometheus:
		return metricsprovider.NewPrometheusMetricsProvider(source.Address)
	case domain.SourceTypeInfluxDB:
		return metricsprovider.NewInfluxDBMetricsProvider()
	default:
		wrErr := errors.Wrap(ErrUnexpectedProviderType, string(source.Type))
		panic(wrErr)
	}
}
