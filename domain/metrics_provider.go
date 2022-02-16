package domain

import "time"

type SeriesItem struct {
	Timestamp time.Time
	Value     float64
}

type MetricsProvider interface {
	LoadSeries(selector string, from, till time.Time, step time.Duration) ([]SeriesItem, error)
}
