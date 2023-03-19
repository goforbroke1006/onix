package domain

import (
	"context"
	"time"
)

// MeasurementRow keep metric pair timestamp-value.
type MeasurementRow struct {
	Moment time.Time
	Value  float64
}

type MeasurementRepository interface {
	Store(
		ctx context.Context,
		sourceID string,
		criteriaID int64,
		moment time.Time,
		value float64) error

	StoreBatch(
		ctx context.Context,
		sourceID string,
		criteriaID int64,
		measurements []MeasurementRow) error
}

type MeasurementService interface {
	GetOrPull(
		ctx context.Context,
		source Source,
		criteria Criteria,
		from, till time.Time, step time.Duration,
	) ([]MeasurementRow, error)
}
