package cmd

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	apiSpec "github.com/goforbroke1006/onix/api/stub_prometheus"
	"github.com/goforbroke1006/onix/internal/component/stub/prometheus"
	pkgEcho "github.com/goforbroke1006/onix/pkg/echo"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewStubPrometheusCmd create prometheus stub cobra-command
func NewStubPrometheusCmd() *cobra.Command {
	const (
		baseURL = "api/v1"
	)

	return &cobra.Command{
		Use: "prometheus",
		Run: func(cmd *cobra.Command, args []string) {
			httpAddr := viper.GetString("server.http.stub.prometheus")

			var (
				logger = log.NewLogger()
			)

			router := echo.New()
			router.Use(middleware.CORS())
			router.HTTPErrorHandler = pkgEcho.ErrorHandler(logger)
			server := prometheus.NewServer()

			apiSpec.RegisterHandlersWithBaseURL(router, server, baseURL)
			if err := router.Start(httpAddr); err != http.ErrServerClosed {
				logger.WithErr(err).Fatal("can't run server")
			}
		},
	}
}
