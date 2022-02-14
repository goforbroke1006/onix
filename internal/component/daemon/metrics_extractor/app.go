package metrics_extractor

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/external/prom"
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
		for _, service := range services {
			criteriaList, err := app.criteriaRepo.GetAll(service.Title)
			if err != nil {
				app.logger.WithErr(err).Warn("can't find criteria list for service", service.Title)
				continue
			}

			criteriaWg := sync.WaitGroup{}

			for _, cr := range criteriaList {
				criteriaWg.Add(1)

				go func(cr domain.Criteria) {
					defer func() {
						criteriaWg.Done()
					}()

					client := prom.NewClient(source.Address)

					var (
						startAt = time.Now().Add(period)
						stopAt  = time.Now()
					)

					resp, err := client.QueryRange(cr.Selector, startAt, stopAt, cr.PullPeriod, 10*time.Second)
					if err != nil {
						app.logger.WithErr(err).Warn("can't extract metric", cr.Title, "from", source.Address)
						return
					}

					if len(resp.Data.Result) == 0 {
						fmt.Printf("no '%s' metric for day %s\n", cr.Title, startAt.Format("2006 Jan 02"))
						return
					}

					batch := make([]domain.MeasurementRow, 0, len(resp.Data.Result[0].Values))
					for _, gv := range resp.Data.Result[0].Values {
						unix := int64(gv[0].(float64))
						moment := time.Unix(unix, 0)
						value, _ := strconv.ParseFloat(gv[1].(string), 64)

						batch = append(batch, domain.MeasurementRow{
							Moment: moment,
							Value:  value,
						})
					}
					if err := app.measurementRepo.StoreBatch(source.ID, cr.ID, batch); err != nil {
						app.logger.WithErr(err).Warn("can't extract metric", cr.Title, "from", source.Address)
					}

					app.logger.Infof("extract %s '%s' %d metrics", service.Title, cr.Title, len(batch))
				}(cr)

			}

			criteriaWg.Wait()
		}
	}
}
