package common

import (
	"fmt"

	"github.com/spf13/viper"
)

// GetDbConnString returns db connection string from viper settings
func GetDbConnString() string {
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
