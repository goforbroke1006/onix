package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"

	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/internal/service"
)

func NewUtilLoadHistoricalMetrics() *cobra.Command {
	var (
		serviceName string
		sourceID    int64
		dateFrom    string
		dateTill    string
	)
	cmd := &cobra.Command{
		Use: "load-historical-metrics",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := pgxpool.Connect(context.Background(), common.GetDbConnString())
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			var (
				sourceRepo      = repository.NewSourceRepository(conn)
				criteriaRepo    = repository.NewCriteriaRepository(conn)
				measurementRepo = repository.NewMeasurementRepository(conn)
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
			till = till.Add(24 * time.Hour)

			var (
				startAt = from
				stopAt  = from.Add(24 * time.Hour).Add(-1 * time.Second)
			)

			for {
				if stopAt.After(till) {
					break
				}

				fmt.Printf("loading %s (%d criteria)\n",
					startAt.Format("2006 Jan 02"), len(criteriaList))

				for _, cr := range criteriaList {
					source, err := sourceRepo.Get(sourceID)
					if err != nil {
						panic(err)
					}
					provider := service.NewMetricsProvider(*source)

					series, err := provider.LoadSeries(cr.Selector, startAt, stopAt, time.Duration(cr.GroupingInterval))
					if err != nil {
						panic(err)
					}

					fmt.Printf("for criteria '%s' loaded series len=%d\n",
						cr.Title, len(criteriaList))
					if len(series) == 0 {
						fmt.Printf("no '%s' metric for day %s\n", cr.Title, startAt.Format("2006 Jan 02"))
						continue
					}

					fmt.Printf("load '%s' metric for day %s\n", cr.Title, startAt.Format("2006 Jan 02"))

					batch := make([]domain.MeasurementRow, 0, len(series))
					for _, item := range series {
						batch = append(batch, domain.MeasurementRow{
							Moment: item.Timestamp,
							Value:  item.Value,
						})
					}
					if err := measurementRepo.StoreBatch(source.ID, cr.ID, batch); err != nil {
						panic(err)
					}
				}

				startAt = startAt.Add(24 * time.Hour)
				stopAt = stopAt.Add(24 * time.Hour)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&serviceName, "service", "", "Service name")
	cmd.PersistentFlags().Int64Var(&sourceID, "source", 0, "Source ID from what need to pull data")
	cmd.PersistentFlags().StringVar(&dateFrom, "from", "", "Time range start")
	cmd.PersistentFlags().StringVar(&dateTill, "till", "", "Time range stop")
	return cmd
}
