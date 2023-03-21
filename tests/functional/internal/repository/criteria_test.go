package repository //nolint:testpackage

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/internal/common"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/tests"
)

func TestGetAll(t *testing.T) { //nolint:paralleltest
	connString := common.GetTestConnectionStrings()

	ctx := context.Background()

	db, dbErr := sqlx.ConnectContext(ctx, "postgres", connString)
	if dbErr != nil {
		t.Skip(dbErr)
	}
	defer func() { _ = db.Close() }()

	if err := tests.LoadFixture(db, "./criteria_test.fixture.sql"); err != nil {
		t.Fatal(err)
	}

	criteriaRepository := repository.NewCriteriaRepository(db)

	criteriaList, critListErr := criteriaRepository.GetAll(ctx, "foo/bar/backend")
	assert.Nil(t, critListErr)
	assert.Equal(t, 3, len(criteriaList))
}
