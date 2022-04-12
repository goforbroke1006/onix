package service

import (
	"context"
	"time"

	"github.com/goforbroke1006/onix/domain"
)

type MeasurementService interface {
	GetOrPull(
		ctx context.Context,
		source domain.Source,
		criteria domain.Criteria,
		from, till time.Time, step time.Duration,
	) ([]domain.MeasurementRow, error)
}
