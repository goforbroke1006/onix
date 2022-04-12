package repository

import (
	"time"

	"github.com/goforbroke1006/onix/domain"
)

// ServiceRepository describes how to manage Service in db.
type ServiceRepository interface {
	Store(title string) error
	GetAll() ([]domain.Service, error)
}

// ReleaseRepository describes how to manage Release in db.
type ReleaseRepository interface {
	Store(serviceName string, releaseName string, startAt time.Time) error
	GetReleases(serviceName string, from, till time.Time) ([]domain.Release, error)
	GetByName(serviceName, releaseName string) (*domain.Release, error)
	GetNextAfter(serviceName, releaseName string) (*domain.Release, error)
	GetLast(serviceName string) (*domain.Release, error)
	GetNLasts(serviceName string, count uint) ([]domain.Release, error)
	GetAll(serviceName string) ([]domain.Release, error)
}

type SourceRepository interface {
	Create(title string, kind domain.SourceType, address string) (int64, error)
	Get(identifier int64) (*domain.Source, error)
	GetAll() ([]domain.Source, error)
}

// CriteriaRepository describe methods for managing Criteria in db.
type CriteriaRepository interface {
	Create(
		serviceName, title string,
		selector string,
		expectedDir domain.DynamicDirType,
		interval domain.GroupingIntervalType,
	) (int64, error)
	GetAll(serviceName string) ([]domain.Criteria, error)
	GetByID(identifier int64) (domain.Criteria, error)
}

// MeasurementRepository describes methods how manage MeasurementRow in db.
type MeasurementRepository interface {
	Store(sourceID, criteriaID int64, moment time.Time, value float64) error
	StoreBatch(sourceID, criteriaID int64, measurements []domain.MeasurementRow) error
	GetBy(sourceID, criteriaID int64, from, till time.Time) ([]domain.MeasurementRow, error)
	Count(sourceID, criteriaID int64, from, till time.Time) (int64, error)
	GetForPoints(sourceID, criteriaID int64, points []time.Time) ([]domain.MeasurementRow, error)
}
