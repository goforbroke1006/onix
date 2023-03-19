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

	"github.com/goforbroke1006/onix/internal/component/stub/prometheus/impl"
	"github.com/goforbroke1006/onix/internal/component/stub/prometheus/spec"
	pkgEcho "github.com/goforbroke1006/onix/pkg/echo"
)

// NewStubPrometheusCmd create prometheus stub cobra-command.
func NewStubPrometheusCmd() *cobra.Command {
	return &cobra.Command{ //nolint:exhaustivestruct
		Use: "prometheus",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			router := echo.New()
			router.Use(middleware.CORS())
			router.HTTPErrorHandler = pkgEcho.ErrorHandler()
			server := impl.NewHandlers()
			spec.RegisterHandlers(router, server)
			go func() {
				if err := router.Start("0.0.0.0:8080"); errors.Is(err, http.ErrServerClosed) {
					zap.L().Fatal("can't run server", zap.Error(err))
				}
			}()

			<-ctx.Done()
		},
	}
}
