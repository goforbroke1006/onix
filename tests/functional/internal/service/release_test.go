//go:build functional
// +build functional

package service

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/cmd"
	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/internal/service"
)

func TestGetReleases(t *testing.T) {
	_ = cmd.ExecuteCmdTree()

	connString := common.GetDbConnString()
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fixtureData, err := ioutil.ReadFile("./release_test.fixture.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = conn.Exec(context.TODO(), string(fixtureData)); err != nil {
		t.Fatal(err)
	}

	var releaseRepository = repository.NewReleaseRepository(conn)
	releaseService := service.NewReleaseService(releaseRepository)

	from, _ := time.Parse("2006-01-02 15:04:05", "2020-10-25 00:00:00")
	till, _ := time.Parse("2006-01-02 15:04:05", "2020-11-06 00:00:00")
	ranges, err := releaseService.GetReleases("foo/bar/backend", from, till)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, len(ranges))
	assert.Equal(t, time.Date(2020, time.October, 25, 0, 0, 0, 0, time.UTC), ranges[0].StartAt)
	assert.Equal(t, time.Date(2020, time.October, 25, 23, 59, 59, 0, time.UTC), ranges[0].StopAt)

	assert.Equal(t, time.Date(2020, time.November, 6, 0, 0, 0, 0, time.UTC), ranges[2].StartAt)
	assert.Equal(t, time.Date(2020, time.November, 13, 23, 59, 59, 0, time.UTC), ranges[2].StopAt)
}

func TestGetReleasesInTheEndOfList(t *testing.T) {
	_ = cmd.ExecuteCmdTree()

	connString := common.GetDbConnString()
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fixtureData, err := ioutil.ReadFile("./release_test.fixture.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = conn.Exec(context.TODO(), string(fixtureData)); err != nil {
		t.Fatal(err)
	}

	var releaseRepository = repository.NewReleaseRepository(conn)
	releaseService := service.NewReleaseService(releaseRepository)

	from, _ := time.Parse("2006-01-02 15:04:05", "2020-11-28 00:00:00")
	till, _ := time.Parse("2006-01-02 15:04:05", "2020-12-26 00:00:00")
	ranges, err := releaseService.GetReleases("foo/bar/backend", from, till)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, len(ranges))

	assert.Equal(t, time.Date(2020, time.December, 26, 0, 0, 0, 0, time.UTC), ranges[2].StartAt)
	assert.Less(t, time.Now().UTC().Sub(ranges[2].StopAt), 2*time.Second)
}
