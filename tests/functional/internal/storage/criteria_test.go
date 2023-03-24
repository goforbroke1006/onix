package storage //nolint:testpackage

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/internal/storage"
	"github.com/goforbroke1006/onix/tests"
)

func TestGetAll(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()

	db, dbErr := GetTestDBConn(ctx)
	if dbErr != nil {
		t.Skip(dbErr)
	}
	defer func() { _ = db.Close() }()

	if err := tests.LoadFixture(db, "./criteria_test.fixture.sql"); err != nil {
		t.Fatal(err)
	}

	criteriaRepository := storage.NewCriteriaStorage(db)

	criteriaList, critListErr := criteriaRepository.GetAll(ctx, "foo/bar/backend")
	assert.Nil(t, critListErr)
	assert.Equal(t, 3, len(criteriaList))
}
