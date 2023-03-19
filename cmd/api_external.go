package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/goforbroke1006/onix/internal/common"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/impl"
	apiSpec "github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
	"github.com/goforbroke1006/onix/internal/repository"
	pkgEcho "github.com/goforbroke1006/onix/pkg/echo"
)

// NewAPIExternalCmd create new system backend cobra-command.
func NewAPIExternalCmd() *cobra.Command {
	return &cobra.Command{ // nolint:exhaustivestruct
		Use: "external",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			httpAddr := viper.GetString("handlers.http.api.system")

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
			)

			router := echo.New()
			router.Use(middleware.CORS())
			router.HTTPErrorHandler = pkgEcho.ErrorHandler()
			handlers := impl.NewHandlers(serviceRepo, sourceRepo, criteriaRepo, releaseRepo)
			apiSpec.RegisterHandlers(router, handlers)
			go func() {
				if startErr := router.Start(httpAddr); errors.Is(err, http.ErrServerClosed) {
					zap.L().Fatal("start server fail", zap.Error(startErr))
				}
			}()

			<-ctx.Done()
		},
	}
}