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

func MustParsePullPeriodType(s string) PullPeriodType {
	switch s {
	case "30s":
		return PullPeriodType30s
	case "1m":
		return PullPeriodType1m
	case "2m":
		return PullPeriodType2m
	case "5m":
		return PullPeriodType5m
	case "15":
		return PullPeriodType15m
	}

	panic(fmt.Errorf("unexpected pull period %s", s))
}

type PullPeriodType time.Duration

const (
	PullPeriodType30s = PullPeriodType(30 * time.Second)
	PullPeriodType1m  = PullPeriodType(1 * time.Minute)
	PullPeriodType2m  = PullPeriodType(2 * time.Minute)
	PullPeriodType5m  = PullPeriodType(5 * time.Minute)
	PullPeriodType15m = PullPeriodType(15 * time.Minute)
)

func (ppt PullPeriodType) String() string {
	switch ppt {
	case PullPeriodType30s:
		return "30s"
	case PullPeriodType1m:
		return "1m"
	case PullPeriodType2m:
		return "2m"
	case PullPeriodType5m:
		return "5m"
	case PullPeriodType15m:
		return "15m"
	}
	return ""
}

type Criteria struct {
	ID          int64
	Service     string
	Title       string
	Selector    string
	ExpectedDir DynamicDirType
	PullPeriod  time.Duration
}

type CriteriaRepository interface {
	Create(
		serviceName, title string,
		selector string,
		expectedDir DynamicDirType,
		pullPeriod PullPeriodType,
	) (int64, error)

	GetAll(serviceName string) ([]Criteria, error)
}
