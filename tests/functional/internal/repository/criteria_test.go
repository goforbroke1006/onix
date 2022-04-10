package repository // nolint:testpackage

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/tests"
)

func TestGetAll(t *testing.T) { // nolint:paralleltest
	connString := common.GetTestConnectionStrings()

	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		t.Skip(err)
	}
	defer conn.Close()

	if err := tests.LoadFixture(conn, "./criteria_test.fixture.sql"); err != nil {
		t.Fatal(err)
	}

	criteriaRepository := repository.NewCriteriaRepository(conn)

	criteriaList, err := criteriaRepository.GetAll("foo/bar/backend")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, len(criteriaList))
}
