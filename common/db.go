package common

import (
	"fmt"

	"github.com/spf13/viper"
)

// GetDBConnString returns db connection string from viper settings.
func GetDBConnString() string {
	var (
		user   = viper.GetString("db.user")
		pass   = viper.GetString("db.pass")
		host   = viper.GetString("db.host")
		port   = viper.GetInt("db.port")
		dbname = viper.GetString("db.dbname")
	)

	target := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, dbname)

	return target
}

func GetTestConnectionStrings() string {
	const (
		user   = "onix"
		pass   = "onix"
		host   = "127.0.0.1"
		port   = 5432
		dbname = "onix"
	)

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, dbname)
}
