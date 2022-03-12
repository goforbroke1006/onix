package service

import (
	"time"

	"github.com/goforbroke1006/onix/domain"
)

// NewReleaseService creates service for manipulate with release data
func NewReleaseService(repo domain.ReleaseRepository) *releaseService {
	return &releaseService{
		repo: repo,
	}
}

var (
	_ domain.ReleaseService = &releaseService{}
)

type releaseService struct {
	repo domain.ReleaseRepository
}

func (svc releaseService) GetReleases(serviceName string, from, till time.Time) ([]domain.ReleaseTimeRange, error) {
	releases, err := svc.repo.GetReleases(serviceName, from, till)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, nil
	}

	ranges := make([]domain.ReleaseTimeRange, 0, len(releases))

	for i := 0; i <= len(releases)-2; i++ {
		ranges = append(ranges, domain.ReleaseTimeRange{
			ID:      releases[i].ID,
			Service: releases[i].Service,
			Name:    releases[i].Name,
			StartAt: releases[i].StartAt,
			StopAt:  releases[i+1].StartAt.Add(-1 * time.Second),
		})
	}

	lastIndex := len(releases) - 1
	ranges = append(ranges, domain.ReleaseTimeRange{
		ID:      releases[lastIndex].ID,
		Service: releases[lastIndex].Service,
		Name:    releases[lastIndex].Name,
		StartAt: releases[lastIndex].StartAt,
		StopAt:  time.Now(),
	})

	afterLast, err := svc.repo.GetNextAfter(serviceName, releases[lastIndex].Name)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	if afterLast != nil {
		ranges[len(ranges)-1].StopAt = afterLast.StartAt.Add(-1 * time.Second)
	}

	return ranges, nil
}

func (svc releaseService) GetAll(serviceName string) ([]domain.ReleaseTimeRange, error) {
	releases, err := svc.repo.GetAll(serviceName)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, nil
	}

	ranges := make([]domain.ReleaseTimeRange, 0, len(releases))

	for i := 0; i <= len(releases)-2; i++ {
		ranges = append(ranges, domain.ReleaseTimeRange{
			ID:      releases[i].ID,
			Service: releases[i].Service,
			Name:    releases[i].Name,
			StartAt: releases[i].StartAt,
			StopAt:  releases[i+1].StartAt.Add(-1 * time.Second),
		})
	}

	lastIndex := len(releases) - 1
	ranges = append(ranges, domain.ReleaseTimeRange{
		ID:      releases[lastIndex].ID,
		Service: releases[lastIndex].Service,
		Name:    releases[lastIndex].Name,
		StartAt: releases[lastIndex].StartAt,
		StopAt:  time.Now().Truncate(time.Second).UTC(),
	})

	return ranges, nil
}

func (svc releaseService) GetByName(serviceName, releaseName string) (*domain.ReleaseTimeRange, error) {
	current, err := svc.repo.GetByName(serviceName, releaseName)
	if err != nil {
		return nil, err
	}

	next, err := svc.repo.GetNextAfter(serviceName, releaseName)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}

	timeRange := domain.ReleaseTimeRange{
		ID:      current.ID,
		Service: serviceName,
		Name:    releaseName,
		StartAt: current.StartAt,
	}

	if next != nil {
		timeRange.StopAt = next.StartAt.Add(-1 * time.Second)
	} else {
		timeRange.StopAt = time.Now().UTC()
	}

	return &timeRange, nil
}
