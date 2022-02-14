package domain

import (
	"time"
)

type Release struct {
	ID      int64
	Service string
	Name    string
	StartAt time.Time
}

type ReleaseTimeRange struct {
	ID      int64
	Service string
	Name    string
	StartAt time.Time
	StopAt  time.Time
}

type ReleaseRepository interface {
	Store(serviceName string, releaseName string, startAt time.Time) error
	GetReleases(serviceName string, from, till time.Time) ([]Release, error)
	GetByName(serviceName, releaseName string) (*Release, error)
	GetNextAfter(serviceName, releaseName string) (*Release, error)
	GetLast(serviceName string) (*Release, error)
	GetAll(serviceName string) ([]Release, error)
}

type ReleaseService interface {
	GetAll(serviceName string) ([]ReleaseTimeRange, error)
	GetReleases(serviceName string, from, till time.Time) ([]ReleaseTimeRange, error)
	GetByName(serviceName, releaseName string) (*ReleaseTimeRange, error)
}
