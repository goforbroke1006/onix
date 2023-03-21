package tests

import (
	"context"
	"os"

	"github.com/jmoiron/sqlx"
)

func LoadFixture(db *sqlx.DB, filename string) error {
	fixtureData, loadFileErr := os.ReadFile(filename)
	if loadFileErr != nil {
		return loadFileErr //nolint:wrapcheck
	}

	if _, execErr := db.ExecContext(context.TODO(), string(fixtureData)); execErr != nil {
		return execErr //nolint:wrapcheck
	}

	return nil
}
