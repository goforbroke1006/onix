package cmd

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-main"
	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/component/api/dashboardmain"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/internal/service"
	pkgEcho "github.com/goforbroke1006/onix/pkg/echo"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewAPIDashboardMainCmd creates dashboard-main backend cobra-command.
func NewAPIDashboardMainCmd() *cobra.Command {
	const (
		baseURL = "api/dashboard-main"
	)

	return &cobra.Command{ // nolint:exhaustivestruct
		Use: "dashboard-main",
		Run: func(cmd *cobra.Command, args []string) {
			httpAddr := viper.GetString("server.http.api.dashboard_main")

			connString := common.GetDBConnString()
			conn, err := pgxpool.Connect(context.Background(), connString)
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			var (
				serviceRepo    = repository.NewServiceRepository(conn)
				releaseSvc     = service.NewReleaseService(repository.NewReleaseRepository(conn))
				sourceRepo     = repository.NewSourceRepository(conn)
				criteriaRepo   = repository.NewCriteriaRepository(conn)
				measurementSvc = service.NewMeasurementService(repository.NewMeasurementRepository(conn))
				logger         = log.NewLogger()
			)

			router := echo.New()
			router.Use(middleware.CORS())
			router.HTTPErrorHandler = pkgEcho.ErrorHandler(logger)
			server := dashboardmain.NewHandlers(serviceRepo, releaseSvc, sourceRepo, criteriaRepo, measurementSvc, logger)

			apiSpec.RegisterHandlersWithBaseURL(router, server, baseURL)
			if err := router.Start(httpAddr); errors.Is(err, http.ErrServerClosed) {
				logger.WithErr(err).Fatal("can't run server")
			}
		},
	}
}
