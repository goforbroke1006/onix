package domain

import (
	"time"
)

// Release keeps data about service's release.
type Release struct {
	Service string
	Tag     string
	StartAt time.Time
}

// ReleaseTimeRange looks like Release but has StopAt time, that calculated from db data.
type ReleaseTimeRange struct {
	Service string
	Tag     string
	StartAt time.Time
	StopAt  time.Time
}

// ReleaseStorage describes how to manage Release in db.
type ReleaseStorage interface {
	Store(serviceName string, releaseName string, startAt time.Time) error
	GetReleases(serviceName string, from, till time.Time) ([]Release, error)
	GetByName(serviceName, tagName string) (*Release, error)
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
