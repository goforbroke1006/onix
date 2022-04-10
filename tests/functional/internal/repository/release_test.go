package repository

import (
	"context"
	"github.com/goforbroke1006/onix/common"
	"io/ioutil"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/internal/repository"
)

func TestGetLast(t *testing.T) {
	connString := common.GetTestConnectionStrings()
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		t.Skip(err)
	}
	defer conn.Close()

	fixtureData, err := ioutil.ReadFile("./release_test.fixture.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = conn.Exec(context.TODO(), string(fixtureData)); err != nil {
		t.Fatal(err)
	}

	releaseRepository := repository.NewReleaseRepository(conn)
	release, err := releaseRepository.GetLast("foo/bar/backend")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "2.1.0", release.Name)
}

func TestGetReleases(t *testing.T) {
	connString := common.GetTestConnectionStrings()
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		t.Skip(err)
	}
	defer conn.Close()

	fixtureData, err := ioutil.ReadFile("./release_test.fixture.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = conn.Exec(context.TODO(), string(fixtureData)); err != nil {
		t.Fatal(err)
	}

	releaseRepository := repository.NewReleaseRepository(conn)

	from, _ := time.Parse("2006-01-02 15:04:05", "2020-10-25 00:00:00")
	till, _ := time.Parse("2006-01-02 15:04:05", "2020-11-06 00:00:00")
	ranges, err := releaseRepository.GetReleases("foo/bar/backend", from, till)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, len(ranges))
	assert.Equal(t, "1.0.0", ranges[0].Name)
	assert.Equal(t, "1.0.1", ranges[1].Name)
	assert.Equal(t, "1.1.0", ranges[2].Name)
}
