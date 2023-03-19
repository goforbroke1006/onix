package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/onix/internal/common"
	"github.com/goforbroke1006/onix/internal/component/daemon/metrics_extractor"
	"github.com/goforbroke1006/onix/internal/repository"
)

// NewDaemonMetricsExtractorCmd create metrics extractor cobra-command.
func NewDaemonMetricsExtractorCmd() *cobra.Command {
	return &cobra.Command{
		Use: "metrics-extractor",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			conn, connErr := pgxpool.Connect(context.Background(), common.GetDBConnString())
			if connErr != nil {
				zap.L().Fatal("db connection fail", zap.Error(connErr))
			}
			defer conn.Close()

			var (
				serviceRepo           = repository.NewServiceRepository(conn)
				criteriaRepository    = repository.NewCriteriaRepository(conn)
				sourceRepository      = repository.NewSourceRepository(conn)
				measurementRepository = repository.NewMeasurementRepository(conn)
			)

			application := metrics_extractor.NewApplication(
				serviceRepo, criteriaRepository, sourceRepository, measurementRepository)
			go func() {
				if runErr := application.Run(ctx); runErr != nil {
					zap.L().Fatal("application stop with fail", zap.Error(runErr))
				}
			}()

			<-ctx.Done()
		},
	}
}

func init() {
	daemonCmd.AddCommand(NewDaemonMetricsExtractorCmd())
}
