package metrics_extractor //nolint:revive,stylecheck

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/service"
)

// NewApplication create metrics extractor app instance.
func NewApplication(
	serviceRepo domain.ServiceStorage,
	criteriaRepo domain.CriteriaStorage,
	sourceRepo domain.SourceStorage,
	measurementRepo domain.MeasurementStorage,
) *Application {
	return &Application{
		serviceRepo:     serviceRepo,
		criteriaRepo:    criteriaRepo,
		sourceRepo:      sourceRepo,
		measurementRepo: measurementRepo,
	}
}

type Application struct {
	serviceRepo     domain.ServiceStorage
	criteriaRepo    domain.CriteriaStorage
	sourceRepo      domain.SourceStorage
	measurementRepo domain.MeasurementStorage
}

func (app Application) Run(ctx context.Context) error {
	const (
		retrievePeriod = 15 * time.Minute

		initialRange = -12 * time.Hour
		anotherRange = -1 * time.Hour
	)

	app.extractMetrics(ctx, initialRange)

	ticker := time.NewTicker(retrievePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			app.extractMetrics(ctx, anotherRange)
		}
	}
}

func (app Application) extractMetrics(ctx context.Context, period time.Duration) { //nolint:funlen,gocognit,cyclop
	sources, sourcesErr := app.sourceRepo.GetAll(ctx)
	if sourcesErr != nil {
		zap.L().Error("can't find sources", zap.Error(sourcesErr))
		return
	}

	services, servicesErr := app.serviceRepo.GetAll(ctx)
	if servicesErr != nil {
		zap.L().Error("can't find services", zap.Error(servicesErr))
		return
	}

	for _, source := range sources {
		for _, svc := range services {
			criteriaList, critListErr := app.criteriaRepo.GetAll(ctx, svc.ID)
			if critListErr != nil {
				zap.L().Error("can't find criteria list for service", zap.Error(critListErr))
				continue
			}

			criteriaWg := sync.WaitGroup{}

			for _, criteria := range criteriaList {
				criteriaWg.Add(1)

				go func(criteria domain.Criteria) {
					defer func() {
						criteriaWg.Done()
					}()

					provider := service.NewMetricsProvider(source)

					var (
						startAt = time.Now().Add(period)
						stopAt  = time.Now()
					)

					resp, err := provider.LoadSeries(ctx,
						criteria.Selector, startAt, stopAt, criteria.Interval)
					if err != nil {
						zap.L().Error("can't extract metric", zap.Error(err))
						return
					}

					if len(resp) == 0 {
						zap.L().Error("no metrics", zap.String("criteria", criteria.Title))
						return
					}

					batch := make([]domain.MeasurementRow, 0, len(resp))
					for _, item := range resp {
						batch = append(batch, domain.MeasurementRow{
							Moment: item.Timestamp,
							Value:  item.Value,
						})
					}

					if err := app.measurementRepo.StoreBatch(ctx, source.ID, criteria.ID, batch); err != nil {
						zap.L().Error("", zap.Error(err))
						return
					}

					zap.L().Info("extract metrics",
						zap.String("service", svc.ID),
						zap.String("criteria", criteria.Title),
						zap.Int("count", len(batch)))
				}(criteria)
			}

			criteriaWg.Wait()
		}
	}
}
