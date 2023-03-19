package service

import (
	"time"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

// NewReleaseService creates service for manipulate with release data.
func NewReleaseService(repo domain.ReleaseRepository) domain.ReleaseService { //nolint:revive,golint
	return &releaseService{repo: repo}
}

var _ domain.ReleaseService = (*releaseService)(nil)

type releaseService struct {
	repo domain.ReleaseRepository
}

func (svc releaseService) GetReleases(serviceName string, from, till time.Time) ([]domain.ReleaseTimeRange, error) {
	releases, err := svc.repo.GetReleases(serviceName, from, till)
	if err != nil {
		return nil, errors.Wrap(err, "can't get releases")
	}

	if len(releases) == 0 {
		return nil, nil
	}

	ranges := make([]domain.ReleaseTimeRange, 0, len(releases))

	for releaseIndex := 0; releaseIndex <= len(releases)-2; releaseIndex++ {
		ranges = append(ranges, domain.ReleaseTimeRange{
			Service: releases[releaseIndex].Service,
			Tag:     releases[releaseIndex].Tag,
			StartAt: releases[releaseIndex].StartAt,
			StopAt:  releases[releaseIndex+1].StartAt.Add(-1 * time.Second),
		})
	}

	lastIndex := len(releases) - 1

	ranges = append(ranges, domain.ReleaseTimeRange{
		Service: releases[lastIndex].Service,
		Tag:     releases[lastIndex].Tag,
		StartAt: releases[lastIndex].StartAt,
		StopAt:  time.Now(),
	})

	afterLast, err := svc.repo.GetNextAfter(serviceName, releases[lastIndex].Tag)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, errors.Wrap(err, "can't get next release")
	}

	if afterLast != nil {
		ranges[len(ranges)-1].StopAt = afterLast.StartAt.Add(-1 * time.Second)
	}

	return ranges, nil
}

func (svc releaseService) GetAll(serviceName string) ([]domain.ReleaseTimeRange, error) {
	releases, err := svc.repo.GetAll(serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "can't get releases")
	}

	if len(releases) == 0 {
		return nil, nil
	}

	ranges := make([]domain.ReleaseTimeRange, 0, len(releases))

	for releaseIndex := 0; releaseIndex <= len(releases)-2; releaseIndex++ {
		ranges = append(ranges, domain.ReleaseTimeRange{
			Service: releases[releaseIndex].Service,
			Tag:     releases[releaseIndex].Tag,
			StartAt: releases[releaseIndex].StartAt,
			StopAt:  releases[releaseIndex+1].StartAt.Add(-1 * time.Second),
		})
	}

	lastIndex := len(releases) - 1

	ranges = append(ranges, domain.ReleaseTimeRange{
		Service: releases[lastIndex].Service,
		Tag:     releases[lastIndex].Tag,
		StartAt: releases[lastIndex].StartAt,
		StopAt:  time.Now().Truncate(time.Second).UTC(),
	})

	return ranges, nil
}

func (svc releaseService) GetByName(serviceName, releaseName string) (*domain.ReleaseTimeRange, error) {
	current, err := svc.repo.GetByName(serviceName, releaseName)
	if err != nil {
		return nil, errors.Wrap(err, "can't get release")
	}

	if current == nil {
		return nil, domain.ErrNotFound
	}

	next, err := svc.repo.GetNextAfter(serviceName, releaseName)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, errors.Wrap(err, "can't get next release")
	}

	var stopAt time.Time
	if next != nil {
		stopAt = next.StartAt.Add(-1 * time.Second)
	} else {
		stopAt = time.Now().UTC()
	}

	timeRange := domain.ReleaseTimeRange{
		Service: serviceName,
		Tag:     releaseName,
		StartAt: current.StartAt,
		StopAt:  stopAt,
	}

	return &timeRange, nil
}
