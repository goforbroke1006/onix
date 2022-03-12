package domain

import (
	"fmt"
	"time"
)

// DynamicDirType define expected progress/regress of metric
type DynamicDirType string

const (
	// DynamicDirTypeIncrease means metric should rise
	DynamicDirTypeIncrease = DynamicDirType("increase")

	// DynamicDirTypeDecrease means metric should fall
	DynamicDirTypeDecrease = DynamicDirType("decrease")

	// DynamicDirTypeEqual means metric should not change
	DynamicDirTypeEqual = DynamicDirType("equal")
)

// MustParseGroupingIntervalType converts string to interval
func MustParseGroupingIntervalType(s string) GroupingIntervalType {
	switch s {
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

	panic(fmt.Errorf("unexpected interval %s", s))
}

// GroupingIntervalType define how often onix should pull data from external providers (prom or influx)
type GroupingIntervalType time.Duration

const (
	// GroupingIntervalType30s means pull every 30 seconds
	GroupingIntervalType30s = GroupingIntervalType(30 * time.Second)

	// GroupingIntervalType1m means pull every minute
	GroupingIntervalType1m = GroupingIntervalType(1 * time.Minute)

	// GroupingIntervalType2m means pull every 2 minutes
	GroupingIntervalType2m = GroupingIntervalType(2 * time.Minute)

	// GroupingIntervalType5m means pull every 5 minutes
	GroupingIntervalType5m = GroupingIntervalType(5 * time.Minute)

	// GroupingIntervalType15m means pull every 15 minutes
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

// Criteria keep data how to evaluate and extract metric from external provider
type Criteria struct {
	ID               int64
	Service          string
	Title            string
	Selector         string
	ExpectedDir      DynamicDirType
	GroupingInterval GroupingIntervalType
}

// CriteriaRepository describe methods for managing Criteria in db
type CriteriaRepository interface {
	Create(
		serviceName, title string,
		selector string,
		expectedDir DynamicDirType,
		interval GroupingIntervalType,
	) (int64, error)

	GetAll(serviceName string) ([]Criteria, error)
}
