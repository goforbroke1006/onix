package metrics_extractor

import (
	"fmt"
	"sync"
	"time"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/service"
	"github.com/goforbroke1006/onix/pkg/log"
)

func NewApplication(
	serviceRepo domain.ServiceRepository,
	criteriaRepo domain.CriteriaRepository,
	sourceRepo domain.SourceRepository,
	measurementRepo domain.MeasurementRepository,

	logger log.Logger,
) *application {
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

func (app application) Run() {
	app.extractMetrics(-12 * time.Hour)

	ticker := time.NewTicker(15 * time.Minute)

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

func (app application) extractMetrics(period time.Duration) {
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

			for _, cr := range criteriaList {
				criteriaWg.Add(1)

				go func(cr domain.Criteria) {
					defer func() {
						criteriaWg.Done()
					}()

					provider := service.NewMetricsProvider(source)

					var (
						startAt = time.Now().Add(period)
						stopAt  = time.Now()
					)

					resp, err := provider.LoadSeries(cr.Selector, startAt, stopAt, cr.GroupingInterval)
					if err != nil {
						app.logger.WithErr(err).Warn("can't extract metric", cr.Title, "from", source.Address)
						return
					}

					if len(resp) == 0 {
						fmt.Printf("no '%s' metric for day %s\n", cr.Title, startAt.Format("2006 Jan 02"))
						return
					}

					batch := make([]domain.MeasurementRow, 0, len(resp))
					for _, item := range resp {
						batch = append(batch, domain.MeasurementRow{
							Moment: item.Timestamp,
							Value:  item.Value,
						})
					}
					if err := app.measurementRepo.StoreBatch(source.ID, cr.ID, batch); err != nil {
						app.logger.WithErr(err).Warn("can't extract metric", cr.Title, "from", source.Address)
					}

					app.logger.Infof("extract %s '%s' %d metrics", svc.Title, cr.Title, len(batch))
				}(cr)

			}

			criteriaWg.Wait()
		}
	}
}
