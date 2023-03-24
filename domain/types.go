package domain

import "time"

// DirectionType define expected progress/regress of metric.
type DirectionType string

// GroupingIntervalType define how often onix should pull data from external providers (prom or influx).
type GroupingIntervalType = time.Duration
