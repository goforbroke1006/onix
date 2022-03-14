package prom

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const defaultTimeout = 5 * time.Second

// NewClient creates prom API client instance.
func NewClient(addr string) *client { // nolint:revive,golint
	httpClient := http.Client{ // nolint:exhaustivestruct
		Timeout: defaultTimeout,
	}

	return &client{
		httpClient: httpClient,
		addr:       addr,
	}
}

var _ APIClient = &client{} // nolint:exhaustivestruct

type client struct {
	httpClient http.Client
	addr       string
}

func (c client) Query(
	ctx context.Context, query string, timestamp time.Time, timeout time.Duration,
) (*QueryResponse, error) {
	addr := fmt.Sprintf("%s/api/v1/query?query=%s&time=%d&timeout=%d",
		c.addr, url.QueryEscape(query), timestamp.Unix(), int64(timeout.Seconds()))

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	req.Header.Add("Accept", "application/json")

	var (
		response *http.Response
		err      error
	)

	if response, err = http.DefaultClient.Do(req); err != nil {
		return nil, errors.Wrap(err, "can't query API method")
	}
	defer response.Body.Close()

	respBytes, _ := ioutil.ReadAll(response.Body)

	var respObj QueryResponse
	if err = json.Unmarshal(respBytes, &respObj); err != nil {
		return nil, errors.Wrap(err, "can't parse prom api response")
	}

	return &respObj, nil
}

func (c client) QueryRange(
	ctx context.Context, query string, start, end time.Time, step, timeout time.Duration,
) (*QueryRangeResponse, error) {
	addr := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%d&timeout=%d",
		c.addr, url.QueryEscape(query), start.Unix(), end.Unix(), int(step.Seconds()), int(timeout.Seconds()),
	)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	req.Header.Add("Accept", "application/json")

	var (
		response *http.Response
		err      error
	)

	if response, err = http.DefaultClient.Do(req); err != nil {
		return nil, errors.Wrap(err, "can't query-range API method")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Wrap(ErrUnexpectedStatusCode, fmt.Sprintf("%d", response.StatusCode))
	}

	respBytes, _ := ioutil.ReadAll(response.Body)

	var respObj QueryRangeResponse
	if err = json.Unmarshal(respBytes, &respObj); err != nil {
		return nil, errors.Wrap(err, "can't parse prom api response")
	}

	return &respObj, nil
}
