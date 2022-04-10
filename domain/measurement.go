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

// MeasurementRepository describes methods how manage MeasurementRow in db.
type MeasurementRepository interface {
	Store(sourceID, criteriaID int64, moment time.Time, value float64) error
	StoreBatch(sourceID, criteriaID int64, measurements []MeasurementRow) error
	GetBy(sourceID, criteriaID int64, from, till time.Time) ([]MeasurementRow, error)
	Count(sourceID, criteriaID int64, from, till time.Time) (int64, error)
	GetForPoints(sourceID, criteriaID int64, points []time.Time) ([]MeasurementRow, error)
}

type MeasurementService interface {
	GetOrPull(
		ctx context.Context,
		source Source,
		criteria Criteria,
		from, till time.Time, step time.Duration,
	) ([]MeasurementRow, error)
}
