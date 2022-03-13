package cmd

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"

	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/component/daemon/metricsextractor"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/pkg/log"
	"github.com/goforbroke1006/onix/pkg/shutdowner"
)

// NewDaemonMetricsExtractorCmd create metrics extractor cobra-command
func NewDaemonMetricsExtractorCmd() *cobra.Command {
	return &cobra.Command{ // nolint:exhaustivestruct
		Use: "metrics-extractor",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := pgxpool.Connect(context.Background(), common.GetDbConnString())
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			var (
				serviceRepo           = repository.NewServiceRepository(conn)
				criteriaRepository    = repository.NewCriteriaRepository(conn)
				sourceRepository      = repository.NewSourceRepository(conn)
				measurementRepository = repository.NewMeasurementRepository(conn)
				logger                = log.NewLogger()
			)

			application := metricsextractor.NewApplication(
				serviceRepo, criteriaRepository, sourceRepository, measurementRepository,
				logger)
			go application.Run()
			defer application.Stop()

			shutdowner.WaitForShutdown()
		},
	}
}
