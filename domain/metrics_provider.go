package domain

import (
	"context"
	"time"
)

// SeriesItem keep metric's pair time-value.
type SeriesItem struct {
	Timestamp time.Time
	Value     float64
}

// MetricsProvider describes methods how pulling data from external sources.
type MetricsProvider interface {
	LoadSeries(ctx context.Context, selector string, from, till time.Time, step time.Duration) ([]SeriesItem, error)
}
