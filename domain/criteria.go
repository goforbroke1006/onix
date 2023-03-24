package domain

import (
	"context"
)

// Criteria keep data how to evaluate and extract metric from external provider.
type Criteria struct {
	ID        int64
	Service   string
	Title     string
	Selector  string
	Direction DirectionType
	Interval  GroupingIntervalType
}

// CriteriaStorage describe methods for managing Criteria in db.
type CriteriaStorage interface {
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
