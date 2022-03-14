package prom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = 5 * time.Second

// NewClient creates prom API client instance
func NewClient(addr string) *client { // nolint:revive,golint
	httpClient := http.Client{ // nolint:exhaustivestruct
		Timeout: defaultTimeout,
	}
	return &client{
		httpClient: httpClient,
		addr:       addr,
	}
}

var (
	_ APIClient = &client{}
)

type client struct {
	httpClient http.Client
	addr       string
}

func (c client) Query(query string, timestamp time.Time, timeout time.Duration) (*QueryResponse, error) {
	addr := fmt.Sprintf("%s/api/v1/query?query=%s&time=%d&timeout=%d",
		c.addr, url.QueryEscape(query), timestamp.Unix(), int64(timeout.Seconds()))

	req, _ := http.NewRequest(http.MethodGet, addr, nil)
	req.Header.Add("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, _ := ioutil.ReadAll(response.Body)

	var respObj QueryResponse
	err = json.Unmarshal(respBytes, &respObj)
	return &respObj, err
}

func (c client) QueryRange(query string, start, end time.Time, step, timeout time.Duration) (*QueryRangeResponse, error) {
	addr := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d&timeout=%d",
		c.addr, url.QueryEscape(query), start.Unix(), end.Unix(), int(step.Seconds()), int(timeout.Seconds()),
	)
	fmt.Println(addr)

	req, _ := http.NewRequest(http.MethodGet, addr, nil)
	req.Header.Add("Accept", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	respBytes, _ := ioutil.ReadAll(response.Body)

	var respObj QueryRangeResponse
	err = json.Unmarshal(respBytes, &respObj)
	return &respObj, err
}
