package metricsextractor

import (
	"context"
	"sync"
	"time"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/service"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewApplication create metrics extractor app instance.
func NewApplication(
	serviceRepo domain.ServiceRepository,
	criteriaRepo domain.CriteriaRepository,
	sourceRepo domain.SourceRepository,
	measurementRepo domain.MeasurementRepository,

	logger log.Logger,
) *application { //nolint:golint,revive
	return &application{
		serviceRepo:     serviceRepo,
		criteriaRepo:    criteriaRepo,
		sourceRepo:      sourceRepo,
		measurementRepo: measurementRepo,
		logger:          logger,

		stopInit: make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

type application struct {
	serviceRepo     domain.ServiceRepository
	criteriaRepo    domain.CriteriaRepository
	sourceRepo      domain.SourceRepository
	measurementRepo domain.MeasurementRepository

	logger log.Logger

	stopInit chan struct{}
	stopDone chan struct{}
}

const defaultInterval = 15 * time.Minute

func (app application) Run() {
	app.extractMetrics(-12 * time.Hour)
	ticker := time.NewTicker(defaultInterval)

Loop:
	for {
		select {
		case <-app.stopInit:
			break Loop
		case <-ticker.C:
			app.extractMetrics(-1 * time.Hour)
		}
	}

	ticker.Stop()
	app.stopDone <- struct{}{}
}

func (app application) Stop() {
	app.stopInit <- struct{}{}
	<-app.stopDone
}

func (app application) extractMetrics(period time.Duration) { // nolint:funlen,gocognit,cyclop
	sources, err := app.sourceRepo.GetAll()
	if err != nil {
		app.logger.WithErr(err).Warn("can't find sources")

		return
	}

	services, err := app.serviceRepo.GetAll()
	if err != nil {
		app.logger.WithErr(err).Warn("can't find services")

		return
	}

	for _, source := range sources {
		for _, svc := range services {
			criteriaList, err := app.criteriaRepo.GetAll(svc.Title)
			if err != nil {
				app.logger.WithErr(err).Warn("can't find criteria list for service", svc.Title)

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

					resp, err := provider.LoadSeries(context.TODO(),
						criteria.Selector, startAt, stopAt, time.Duration(criteria.GroupingInterval))
					if err != nil {
						app.logger.WithErr(err).Warn("can't extract metric",
							criteria.Title, "from", source.Address)

						return
					}

					if len(resp) == 0 {
						app.logger.Infof("no '%s' metric for day %s\n",
							criteria.Title, startAt.Format("2006 Jan 02"))

						return
					}

					batch := make([]domain.MeasurementRow, 0, len(resp))
					for _, item := range resp {
						batch = append(batch, domain.MeasurementRow{
							Moment: item.Timestamp,
							Value:  item.Value,
						})
					}

					if err := app.measurementRepo.StoreBatch(source.ID, criteria.ID, batch); err != nil {
						app.logger.WithErr(err).Warn("can't extract metric", criteria.Title, "from", source.Address)
					}

					app.logger.Infof("extract %s '%s' %d metrics", svc.Title, criteria.Title, len(batch))
				}(criteria)
			}

			criteriaWg.Wait()
		}
	}
}
