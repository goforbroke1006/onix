package domain

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// DirectionType define expected progress/regress of metric.
type DirectionType string

const (
	// DynamicDirTypeIncrease means metric should rise.
	DynamicDirTypeIncrease = DirectionType("increase")

	// DynamicDirTypeDecrease means metric should fall.
	DynamicDirTypeDecrease = DirectionType("decrease")

	// DynamicDirTypeEqual means metric should not change.
	DynamicDirTypeEqual = DirectionType("equal")
)

var ErrParseGroupIntervalFailed = errors.New("parse group interval failed")

// MustParseGroupingIntervalType converts string to interval.
func MustParseGroupingIntervalType(str string) GroupingIntervalType {
	switch str {
	case "30s":
		return GroupingIntervalType30s
	case "1m":
		return GroupingIntervalType1m
	case "2m":
		return GroupingIntervalType2m
	case "5m":
		return GroupingIntervalType5m
	case "15":
		return GroupingIntervalType15m
	}

	panic(errors.Wrap(ErrParseGroupIntervalFailed, str))
}

// GroupingIntervalType define how often onix should pull data from external providers (prom or influx).
type GroupingIntervalType time.Duration

const (
	// GroupingIntervalType30s means pull every 30 seconds.
	GroupingIntervalType30s = GroupingIntervalType(30 * time.Second)

	// GroupingIntervalType1m means pull every minute.
	GroupingIntervalType1m = GroupingIntervalType(1 * time.Minute)

	// GroupingIntervalType2m means pull every 2 minutes.
	GroupingIntervalType2m = GroupingIntervalType(2 * time.Minute)

	// GroupingIntervalType5m means pull every 5 minutes.
	GroupingIntervalType5m = GroupingIntervalType(5 * time.Minute)

	// GroupingIntervalType15m means pull every 15 minutes.
	GroupingIntervalType15m = GroupingIntervalType(15 * time.Minute)
)

func (ppt GroupingIntervalType) String() string {
	switch ppt {
	case GroupingIntervalType30s:
		return "30s"
	case GroupingIntervalType1m:
		return "1m"
	case GroupingIntervalType2m:
		return "2m"
	case GroupingIntervalType5m:
		return "5m"
	case GroupingIntervalType15m:
		return "15m"
	}

	return ""
}

// Criteria keep data how to evaluate and extract metric from external provider.
type Criteria struct {
	ID        int64
	Service   string
	Title     string
	Selector  string
	Direction DirectionType
	Interval  GroupingIntervalType
}

// CriteriaRepository describe methods for managing Criteria in db.
type CriteriaRepository interface {
	Create(
		ctx context.Context,
		serviceName, title string,
		selector string,
		expectedDir DirectionType,
		interval GroupingIntervalType,
	) (int64, error)
	GetAll(ctx context.Context, serviceName string) ([]Criteria, error)
	GetByID(ctx context.Context, identifier int64) (Criteria, error)
}
