package common

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func OpenDBConn(ctx context.Context) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", GetDBConnString())
}

// GetDBConnString returns db connection string from viper settings.
func GetDBConnString() string {
	var (
		host   = viper.GetString("db.host")
		port   = viper.GetInt("db.port")
		user   = viper.GetString("db.user")
		pass   = viper.GetString("db.pass")
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
