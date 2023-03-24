package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/onix/internal/common"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/impl"
	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
	"github.com/goforbroke1006/onix/internal/service"
	"github.com/goforbroke1006/onix/internal/storage"
	pkgEcho "github.com/goforbroke1006/onix/pkg/echo"
)

// NewAPIExternalCmd create new system backend cobra-command.
func NewAPIExternalCmd() *cobra.Command {
	return &cobra.Command{ //nolint:exhaustivestruct
		Use: "external",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			db, dbErr := common.OpenDBConn(ctx)
			if dbErr != nil {
				zap.L().Fatal("db connection fail", zap.Error(dbErr))
			}
			defer func() { _ = db.Close() }()

			var (
				serviceRepo  = storage.NewServiceStorage(db)
				sourceRepo   = storage.NewSourceStorage(db)
				criteriaRepo = storage.NewCriteriaStorage(db)
				releaseRepo  = storage.NewReleaseStorage(db)
				releaseSvc   = service.NewReleaseService(releaseRepo)
			)

			router := echo.New()
			router.Use(middleware.CORS())
			router.HTTPErrorHandler = pkgEcho.ErrorHandler()
			handlers := impl.NewHandlers(serviceRepo, sourceRepo, criteriaRepo, releaseRepo, releaseSvc)
			spec.RegisterHandlers(router, handlers)
			go func() {
				if startErr := router.Start("0.0.0.0:8080"); errors.Is(startErr, http.ErrServerClosed) {
					zap.L().Fatal("start server fail", zap.Error(startErr))
				}
			}()

			<-ctx.Done()
		},
	}
}
