package domain

import "time"

const (
	// DynamicDirTypeIncrease means metric should rise.
	DynamicDirTypeIncrease = DirectionType("increase")

	// DynamicDirTypeDecrease means metric should fall.
	DynamicDirTypeDecrease = DirectionType("decrease")

	// DynamicDirTypeEqual means metric should not change.
	DynamicDirTypeEqual = DirectionType("equal")
)

const (
	// GroupingIntervalType30s means pull every 30 seconds.
	GroupingIntervalType30s = 30 * time.Second

	// GroupingIntervalType1m means pull every minute.
	GroupingIntervalType1m = 1 * time.Minute

	// GroupingIntervalType2m means pull every 2 minutes.
	GroupingIntervalType2m = 2 * time.Minute

	// GroupingIntervalType5m means pull every 5 minutes.
	GroupingIntervalType5m = 5 * time.Minute

	// GroupingIntervalType15m means pull every 15 minutes.
	GroupingIntervalType15m = 15 * time.Minute
)
