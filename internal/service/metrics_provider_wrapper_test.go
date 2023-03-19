package service // nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/domain"
)

func TestNewMetricsProvider(t *testing.T) {
	t.Parallel()

	t.Run("get prometheus provider", func(t *testing.T) {
		t.Parallel()

		source := domain.Source{ // nolint:exhaustivestruct
			Kind: domain.SourceTypePrometheus,
		}
		provider := NewMetricsProvider(source)
		assert.NotNil(t, provider, "NewMetricsProvider(%v)", source)
	})

	t.Run("get influx-db provider", func(t *testing.T) {
		t.Parallel()

		source := domain.Source{ // nolint:exhaustivestruct
			Kind: domain.SourceTypeInfluxDB,
		}
		provider := NewMetricsProvider(source)
		assert.NotNil(t, provider, "NewMetricsProvider(%v)", source)
	})

	t.Run("get unknown provider", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		source := domain.Source{ // nolint:exhaustivestruct
			Kind: domain.SourceType("unknown"),
		}
		_ = NewMetricsProvider(source)
	})
}
