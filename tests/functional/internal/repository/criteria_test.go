package repository

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/repository"
)

func TestGetAll(t *testing.T) {
	connString := common.GetTestConnectionStrings()
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		t.Skip(err)
	}
	defer conn.Close()

	fixtureData, err := ioutil.ReadFile("./criteria_test.fixture.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = conn.Exec(context.TODO(), string(fixtureData)); err != nil {
		t.Fatal(err)
	}

	criteriaRepository := repository.NewCriteriaRepository(conn)

	criteriaList, err := criteriaRepository.GetAll("foo/bar/backend")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, len(criteriaList))
}
