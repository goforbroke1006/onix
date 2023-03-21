package cmd

import (
	"context"
	"os"
	"os/signal"

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

			db, dbErr := common.OpenDBConn(ctx)
			if dbErr != nil {
				zap.L().Fatal("db connection fail", zap.Error(dbErr))
			}
			defer func() { _ = db.Close() }()

			var (
				serviceRepo           = repository.NewServiceRepository(db)
				criteriaRepository    = repository.NewCriteriaRepository(db)
				sourceRepository      = repository.NewSourceRepository(db)
				measurementRepository = repository.NewMeasurementRepository(db)
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
