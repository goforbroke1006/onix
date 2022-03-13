package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ExecuteCmdTree enables default settings for viper and initialize cobra commands tree
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

	return rootCmd.Execute()
}

func setupViper() error {
	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "onix")
	viper.SetDefault("db.pass", "onix")
	viper.SetDefault("db.dbname", "onix")

	viper.SetConfigName("onix")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}
