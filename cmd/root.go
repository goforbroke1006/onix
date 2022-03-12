package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ExecuteCmdTree enables default settings for viper and initialize cobra commands tree
func ExecuteCmdTree() error {
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

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	var (
		rootCmd = &cobra.Command{
			Use:   "onix",
			Short: "Onix is performance comparison service",
			Long:  `Onix collect information about releases and compare metrics between its.`,
		}

		apiCmd = &cobra.Command{
			Use:   "api",
			Short: "API handlers",
		}
		daemonCmd = &cobra.Command{
			Use:   "daemon",
			Short: "Periodic task runners",
		}
		stubCmd = &cobra.Command{
			Use:   "stub",
			Short: "Fake external services",
		}
		utilCmd = &cobra.Command{
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
