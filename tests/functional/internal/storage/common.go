package storage

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/goforbroke1006/onix/internal/common"
)

func GetTestDBConn(ctx context.Context) (*sqlx.DB, error) {
	connString := common.GetTestConnectionStrings()
	db, dbErr := sqlx.ConnectContext(ctx, "postgres", connString)
	return db, dbErr
}
