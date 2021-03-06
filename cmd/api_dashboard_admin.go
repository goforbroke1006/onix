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

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-admin"
	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/component/api/dashboardadmin"
	"github.com/goforbroke1006/onix/internal/repository"
	pkgEcho "github.com/goforbroke1006/onix/pkg/echo"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewAPIDashboardAdminCmd creates dashboard-admin backend cobra-command.
func NewAPIDashboardAdminCmd() *cobra.Command {
	const (
		baseURL = "api/dashboard-admin"
	)

	return &cobra.Command{ // nolint:exhaustivestruct
		Use: "dashboard-admin",
		Run: func(cmd *cobra.Command, args []string) {
			httpAddr := viper.GetString("server.http.api.dashboard_admin")

			connString := common.GetDBConnString()
			conn, err := pgxpool.Connect(context.Background(), connString)
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			var (
				serviceRepo  = repository.NewServiceRepository(conn)
				releaseRepo  = repository.NewReleaseRepository(conn)
				sourceRepo   = repository.NewSourceRepository(conn)
				criteriaRepo = repository.NewCriteriaRepository(conn)
				logger       = log.NewLogger()
			)

			router := echo.New()
			router.Use(middleware.CORS())
			router.HTTPErrorHandler = pkgEcho.ErrorHandler(logger)
			server := dashboardadmin.NewHandlers(serviceRepo, releaseRepo, sourceRepo, criteriaRepo, logger)

			apiSpec.RegisterHandlersWithBaseURL(router, server, baseURL)
			if err := router.Start(httpAddr); errors.Is(err, http.ErrServerClosed) {
				logger.WithErr(err).Fatal("can't run server")
			}
		},
	}
}
