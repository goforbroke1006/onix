package storage //nolint:testpackage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/internal/storage"
	"github.com/goforbroke1006/onix/tests"
)

func TestGetLast(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()

	db, dbErr := GetTestDBConn(ctx)
	if dbErr != nil {
		t.Skip(dbErr)
	}
	defer func() { _ = db.Close() }()

	if err := tests.LoadFixture(db, "./release_test.fixture.sql"); err != nil {
		t.Fatal(err)
	}

	releaseRepository := storage.NewReleaseStorage(db)

	release, err := releaseRepository.GetLast("foo/bar/backend")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "2.1.0", release.Tag)
}

func TestGetReleases(t *testing.T) { //nolint:paralleltest
	ctx := context.Background()

	db, dbErr := GetTestDBConn(ctx)
	if dbErr != nil {
		t.Skip(dbErr)
	}
	defer func() { _ = db.Close() }()

	if err := tests.LoadFixture(db, "./release_test.fixture.sql"); err != nil {
		t.Fatal(err)
	}

	releaseRepository := storage.NewReleaseStorage(db)

	from, _ := time.Parse("2006-01-02 15:04:05", "2020-10-25 00:00:00")
	till, _ := time.Parse("2006-01-02 15:04:05", "2020-11-06 00:00:00")

	ranges, err := releaseRepository.GetReleases("foo/bar/backend", from, till)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, len(ranges))
	assert.Equal(t, "1.0.0", ranges[0].Tag)
	assert.Equal(t, "1.0.1", ranges[1].Tag)
	assert.Equal(t, "1.1.0", ranges[2].Tag)
}
