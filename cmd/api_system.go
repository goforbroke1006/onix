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

	apiSpec "github.com/goforbroke1006/onix/api/system"
	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/component/api/system"
	"github.com/goforbroke1006/onix/internal/repository"
	pkgEcho "github.com/goforbroke1006/onix/pkg/echo"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewAPISystemCmd create new system backend cobra-command.
func NewAPISystemCmd() *cobra.Command {
	const (
		baseURL = "api/system"
	)

	return &cobra.Command{ // nolint:exhaustivestruct
		Use: "system",
		Run: func(cmd *cobra.Command, args []string) {
			httpAddr := viper.GetString("server.http.api.system")

			connString := common.GetDBConnString()
			conn, err := pgxpool.Connect(context.Background(), connString)
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			var (
				serviceRepo  = repository.NewServiceRepository(conn)
				sourceRepo   = repository.NewSourceRepository(conn)
				criteriaRepo = repository.NewCriteriaRepository(conn)
				releaseRepo  = repository.NewReleaseRepository(conn)
				logger       = log.NewLogger()
			)

			router := echo.New()
			router.Use(middleware.CORS())
			router.HTTPErrorHandler = pkgEcho.ErrorHandler(logger)
			server := system.NewHandlers(serviceRepo, sourceRepo, criteriaRepo, releaseRepo, logger)

			apiSpec.RegisterHandlersWithBaseURL(router, server, baseURL)
			if err := router.Start(httpAddr); errors.Is(err, http.ErrServerClosed) {
				logger.WithErr(err).Fatal("can't run server")
			}
		},
	}
}
