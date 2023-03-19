package domain

import (
	"time"
)

// Release keeps data about service's release.
type Release struct {
	ID      int64
	Service string
	Tag     string
	StartAt time.Time
}

// ReleaseTimeRange looks like Release but has StopAt time, that calculated from db data.
type ReleaseTimeRange struct {
	ID      int64
	Service string
	Tag     string
	StartAt time.Time
	StopAt  time.Time
}

// ReleaseRepository describes how to manage Release in db.
type ReleaseRepository interface {
	Store(serviceName string, releaseName string, startAt time.Time) error
	GetReleases(serviceName string, from, till time.Time) ([]Release, error)
	GetByName(serviceName, releaseName string) (*Release, error)
	GetNextAfter(serviceName, releaseName string) (*Release, error)
	GetLast(serviceName string) (*Release, error)
	GetNLasts(serviceName string, count uint) ([]Release, error)
	GetAll(serviceName string) ([]Release, error)
}

// ReleaseService helps build ReleaseTimeRange instance.
type ReleaseService interface {
	GetAll(serviceName string) ([]ReleaseTimeRange, error)
	GetReleases(serviceName string, from, till time.Time) ([]ReleaseTimeRange, error)
	GetByName(serviceName, releaseName string) (*ReleaseTimeRange, error)
}
