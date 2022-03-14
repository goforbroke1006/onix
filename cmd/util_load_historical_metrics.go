package cmd

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"

	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/internal/service"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewUtilLoadHistoricalMetrics creates load history util cobra-command.
func NewUtilLoadHistoricalMetrics() *cobra.Command { // nolint:funlen,gocognit,cyclop
	var (
		serviceName string
		sourceID    int64
		dateFrom    string
		dateTill    string
	)

	const oneDay = 24 * time.Hour

	cmd := &cobra.Command{ // nolint:exhaustivestruct
		Use: "load-historical-metrics",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := pgxpool.Connect(context.Background(), common.GetDBConnString())
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			var (
				sourceRepo      = repository.NewSourceRepository(conn)
				criteriaRepo    = repository.NewCriteriaRepository(conn)
				measurementRepo = repository.NewMeasurementRepository(conn)
				logger          = log.NewLogger()
			)

			criteriaList, err := criteriaRepo.GetAll(serviceName)
			if err != nil {
				panic(err)
			}

			var (
				from time.Time
				till time.Time
			)

			from, err = time.Parse("2006-01-02", dateFrom)
			if err != nil {
				panic(err)
			}
			till, err = time.Parse("2006-01-02", dateTill)
			if err != nil {
				panic(err)
			}

			till = till.Add(oneDay)

			var (
				startAt = from
				stopAt  = from.Add(oneDay).Add(-1 * time.Second)
			)

			for {
				if stopAt.After(till) {
					break
				}

				logger.Infof("loading %s (%d criteria)\n",
					startAt.Format("2006 Jan 02"), len(criteriaList))

				for _, criteria := range criteriaList {
					source, err := sourceRepo.Get(sourceID)
					if err != nil {
						panic(err)
					}
					provider := service.NewMetricsProvider(*source)

					series, err := provider.LoadSeries(context.TODO(),
						criteria.Selector, startAt, stopAt, time.Duration(criteria.GroupingInterval))
					if err != nil {
						panic(err)
					}

					logger.Infof("for criteria '%s' loaded series len=%d\n",
						criteria.Title, len(criteriaList))
					if len(series) == 0 {
						logger.Infof("no '%s' metric for day %s\n", criteria.Title, startAt.Format("2006 Jan 02"))

						continue
					}

					logger.Infof("load '%s' metric for day %s\n", criteria.Title, startAt.Format("2006 Jan 02"))

					batch := make([]domain.MeasurementRow, 0, len(series))
					for _, item := range series {
						batch = append(batch, domain.MeasurementRow{
							Moment: item.Timestamp,
							Value:  item.Value,
						})
					}
					if err := measurementRepo.StoreBatch(source.ID, criteria.ID, batch); err != nil {
						panic(err)
					}
				}

				startAt = startAt.Add(oneDay)
				stopAt = stopAt.Add(oneDay)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&serviceName, "service", "", "Service name")
	cmd.PersistentFlags().Int64Var(&sourceID, "source", 0, "Source ID from what need to pull data")
	cmd.PersistentFlags().StringVar(&dateFrom, "from", "", "Time range start")
	cmd.PersistentFlags().StringVar(&dateTill, "till", "", "Time range stop")

	return cmd
}
