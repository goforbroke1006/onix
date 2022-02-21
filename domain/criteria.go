package domain

import (
	"fmt"
	"time"
)

type DynamicDirType string

const (
	DynamicDirTypeIncrease = DynamicDirType("increase")
	DynamicDirTypeDecrease = DynamicDirType("decrease")
	DynamicDirTypeEqual    = DynamicDirType("equal")
)

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

type GroupingIntervalType time.Duration

const (
	GroupingIntervalType30s = GroupingIntervalType(30 * time.Second)
	GroupingIntervalType1m  = GroupingIntervalType(1 * time.Minute)
	GroupingIntervalType2m  = GroupingIntervalType(2 * time.Minute)
	GroupingIntervalType5m  = GroupingIntervalType(5 * time.Minute)
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

type Criteria struct {
	ID               int64
	Service          string
	Title            string
	Selector         string
	ExpectedDir      DynamicDirType
	GroupingInterval GroupingIntervalType
}

type CriteriaRepository interface {
	Create(
		serviceName, title string,
		selector string,
		expectedDir DynamicDirType,
		interval GroupingIntervalType,
	) (int64, error)

	GetAll(serviceName string) ([]Criteria, error)
}
