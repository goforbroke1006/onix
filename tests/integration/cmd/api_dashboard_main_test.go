//go:build integration
// +build integration

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goforbroke1006/onix/internal/component/api/dashboard_main"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/cmd"
	"github.com/goforbroke1006/onix/common"
)

func TestExtractApiHandler(t *testing.T) {
	const (
		hostDockerComposeApiFrontend = "http://127.0.0.1:8082"
		fakeServiceName              = "foo/bar/backend"
		sourceID                     = 1
	)

	_ = cmd.ExecuteCmdTree()

	connString := common.GetDbConnString()
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fixtureData, err := ioutil.ReadFile("./api_dashboard_main_test.fixture.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = conn.Exec(context.TODO(), string(fixtureData)); err != nil {
		t.Fatal(err)
	}

	// send request
	params := []string{
		fmt.Sprintf("service=%s", fakeServiceName),
		fmt.Sprintf("source_id=%d", sourceID),
		"release_one_title=2.0.0",
		"release_one_start=1642877700",
		"release_two_title=2.1.0",
		"release_two_start=1643894976",
		"period=1h",
	}
	addr := fmt.Sprintf("%s/api/dashboard-main/compare?%s", hostDockerComposeApiFrontend, strings.Join(params, "&"))
	resp, err := http.Get(addr)
	if err != nil {
		t.Fatal(err)
	}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	responseObj := dashboard_main.CompareResponse{}
	if err := json.Unmarshal(respBytes, &responseObj); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fakeServiceName, responseObj.Service)
	assert.Equal(t, 2, len(responseObj.Reports))
}
