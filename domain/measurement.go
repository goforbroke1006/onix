package domain

import "time"

type MeasurementRow struct {
	Moment time.Time
	Value  float64
}

type MeasurementShortRow struct {
	Moment time.Time
	Value  float64
}

type MeasurementRepository interface {
	Store(sourceID, criteriaID int64, moment time.Time, value float64) error
	StoreBatch(sourceID, criteriaID int64, measurements []MeasurementRow) error
	GetBy(sourceID, criteriaID int64, from, till time.Time) ([]MeasurementShortRow, error)
	Count(sourceID, criteriaID int64, from, till time.Time) (int64, error)
}
