package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "onix",
		Short: "Onix is performance comparison service",
		Long:  `Onix collect information about releases and compare metrics between its.`,
	}

	apiCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "api",
		Short: "API handlers",
	}
	daemonCmd = &cobra.Command{Use: "daemon"} //nolint:gochecknoglobals
	stubCmd   = &cobra.Command{               //nolint:gochecknoglobals
		Use:   "stub",
		Short: "Fake external services",
	}
	utilCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "util",
		Short: "CLI utils",
	}
)

// ExecuteCmdTree enables default settings for viper and initialize cobra commands tree.
func ExecuteCmdTree() error {
	if err := setupViper(); err != nil {
		return err
	}

	rootCmd.AddCommand(apiCmd, daemonCmd, stubCmd, utilCmd)

	apiCmd.AddCommand(
		NewAPIExternalCmd(),
	)
	stubCmd.AddCommand(
		NewStubPrometheusCmd(),
	)
	utilCmd.AddCommand(
		NewUtilMigrateCmd(),
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

	return nil
}
