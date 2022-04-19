package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

func NewMeasurementService(
	measurementRepo domain.MeasurementRepository,
) domain.MeasurementService {
	return &measurementService{
		measurementRepo:  measurementRepo,
		createProviderFn: NewMetricsProvider,
	}
}

type measurementService struct {
	measurementRepo  domain.MeasurementRepository
	createProviderFn func(source domain.Source) domain.MetricsProvider
}

func (svc measurementService) GetOrPull(
	ctx context.Context,
	source domain.Source,
	criteria domain.Criteria,
	from, till time.Time, step time.Duration,
) ([]domain.MeasurementRow, error) {
	points := svc.getTimePoints(from, till, step)

	measurementRows, err := svc.measurementRepo.GetForPoints(source.ID, criteria.ID, points)
	if err != nil {
		return nil, errors.Wrap(err, "can't get measurements")
	}

	if len(measurementRows) == len(points) {
		return measurementRows, nil
	}

	provider := svc.createProviderFn(source)

	series, err := provider.LoadSeries(ctx, criteria.Selector, from, till, step)
	if err != nil {
		return nil, errors.Wrap(err, "can't pull series")
	}

	batch := make([]domain.MeasurementRow, 0, len(series))
	for _, item := range series {
		batch = append(batch, domain.MeasurementRow{
			Moment: item.Timestamp,
			Value:  item.Value,
		})
	}

	if err := svc.measurementRepo.StoreBatch(source.ID, criteria.ID, batch); err != nil {
		return nil, errors.Wrap(err, "can't store series")
	}

	batchMap := make(map[time.Time]float64, len(batch))
	for _, p := range series {
		batchMap[p.Timestamp] = p.Value
	}

	result := make([]domain.MeasurementRow, 0, len(points))

	for _, timePoint := range points {
		value := 0.0
		if v, ok := batchMap[timePoint]; ok {
			value = v
		}

		result = append(result, domain.MeasurementRow{Moment: timePoint, Value: value})
	}

	return result, nil
}

func (svc measurementService) getTimePoints(from, till time.Time, step time.Duration) []time.Time {
	var result []time.Time
	for t := from; t.Before(till) || t.Equal(till); t = t.Add(step) {
		result = append(result, t)
	}

	return result
}
