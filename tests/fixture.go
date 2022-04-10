package tests

import (
	"context"
	"io/ioutil"

	"github.com/jackc/pgx/v4/pgxpool"
)

func LoadFixture(conn *pgxpool.Pool, filename string) error {
	fixtureData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err // nolint:wrapcheck
	}

	if _, err := conn.Exec(context.TODO(), string(fixtureData)); err != nil {
		return err // nolint:wrapcheck
	}

	return nil
}
