package cmd

import (
	"context"
	"fmt"
	"github.com/goforbroke1006/onix/domain"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"

	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/external/prom"
	"github.com/goforbroke1006/onix/internal/repository"
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

				for _, cr := range criteriaList {
					source, err := sourceRepo.Get(sourceID)
					if err != nil {
						panic(err)
					}
					p := prom.NewClient(source.Address)

					resp, err := p.QueryRange(cr.Selector, startAt, stopAt, cr.PullPeriod, 10*time.Second)
					if err != nil {
						panic(err)
					}

					if len(resp.Data.Result) == 0 {
						fmt.Printf("no '%s' metric for day %s\n", cr.Title, startAt.Format("2006 Jan 02"))
						continue
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
					if err := measurementRepo.StoreBatch(source.ID, cr.ID, batch); err != nil {
						panic(err)
					}
				}
				fmt.Printf("load %s\n", startAt.Format("2006 Jan 02"))

				startAt = startAt.Add(24 * time.Hour)
				stopAt = stopAt.Add(24 * time.Hour)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&serviceName, "service", "", "")
	cmd.PersistentFlags().Int64Var(&sourceID, "source", 0, "")
	cmd.PersistentFlags().StringVar(&dateFrom, "from", "", "")
	cmd.PersistentFlags().StringVar(&dateTill, "till", "", "")
	return cmd
}
