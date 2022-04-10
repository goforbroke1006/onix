package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ExecuteCmdTree enables default settings for viper and initialize cobra commands tree.
func ExecuteCmdTree() error {
	if err := setupViper(); err != nil {
		return err
	}

	var (
		rootCmd = &cobra.Command{ // nolint:exhaustivestruct
			Use:   "onix",
			Short: "Onix is performance comparison service",
			Long:  `Onix collect information about releases and compare metrics between its.`,
		}

		apiCmd = &cobra.Command{ // nolint:exhaustivestruct
			Use:   "api",
			Short: "API handlers",
		}
		daemonCmd = &cobra.Command{ // nolint:exhaustivestruct
			Use:   "daemon",
			Short: "Periodic task runners",
		}
		stubCmd = &cobra.Command{ // nolint:exhaustivestruct
			Use:   "stub",
			Short: "Fake external services",
		}
		utilCmd = &cobra.Command{ // nolint:exhaustivestruct
			Use:   "util",
			Short: "CLI utils",
		}
	)

	rootCmd.AddCommand(apiCmd, daemonCmd, stubCmd, utilCmd)

	apiCmd.AddCommand(
		NewAPISystemCmd(),
		NewAPIDashboardMainCmd(),
		NewAPIDashboardAdminCmd(),
	)
	daemonCmd.AddCommand(
		NewDaemonMetricsExtractorCmd(),
	)
	stubCmd.AddCommand(
		NewStubPrometheusCmd(),
	)
	utilCmd.AddCommand(
		NewUtilLoadHistoricalMetrics(),
	)

	err := rootCmd.Execute()

	return errors.Wrap(err, "can't run root cmd")
}

func setupViper() error {
	const postgresDefaultPort = 5432

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", postgresDefaultPort)
	viper.SetDefault("db.user", "onix")
	viper.SetDefault("db.pass", "onix")
	viper.SetDefault("db.dbname", "onix")

	viper.SetDefault("server.http.api.dashboard_admin", "0.0.0.0:8080")
	viper.SetDefault("server.http.api.dashboard_main", "0.0.0.0:8080")
	viper.SetDefault("server.http.api.system", "0.0.0.0:8080")

	viper.SetConfigName("onix")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.ReadInConfig()
	//if err := viper.ReadInConfig(); err != nil {
	//	return errors.Wrap(err, "can't find config file")
	//}

	return nil
}
