package prometheus // nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	apiSpec "github.com/goforbroke1006/onix/api/stub_prometheus"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	h := NewHandlers()
	assert.NotNil(t, h)
}

func Test_handlers_GetQuery(t *testing.T) {
	t.Parallel()

	var target handlers

	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
	recorder := httptest.NewRecorder()
	echoContext := echo.New().NewContext(req, recorder)

	err := target.GetQuery(echoContext)
	assert.Nil(t, err)
}

func Test_handlers_GetQueryRange(t *testing.T) { // nolint:funlen
	t.Parallel()

	type args struct {
		url    string
		params apiSpec.GetQueryRangeParams
	}

	tests := []struct {
		name         string
		args         args
		wantRespCode int
		wantErr      bool
		wantCount    int
	}{
		{
			name: "negative 1 - empty query",
			args: args{
				url: "https://test.com/api/v1/query_range?query=&start=&end=&spte&",
				params: apiSpec.GetQueryRangeParams{
					Query:   "",
					Start:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
					End:     time.Now().Format(time.RFC3339),
					Step:    "5m",
					Timeout: "5s",
				},
			},
			wantRespCode: http.StatusBadRequest,
			wantErr:      false,
		},
		{
			name: "negative 1 - invalid start",
			args: args{
				url: "",
				params: apiSpec.GetQueryRangeParams{
					Query:   `rate(some_metrics{env="prod"})`,
					Start:   "2022-01-25",
					End:     time.Now().Format(time.RFC3339),
					Step:    "5m",
					Timeout: "5s",
				},
			},
			wantRespCode: http.StatusBadRequest,
			wantErr:      false,
		},
		{
			name: "negative 1 - invalid end",
			args: args{
				url: "",
				params: apiSpec.GetQueryRangeParams{
					Query:   `rate(some_metrics{env="prod"})`,
					Start:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
					End:     "2022-01-25",
					Step:    "5m",
					Timeout: "5s",
				},
			},
			wantRespCode: http.StatusBadRequest,
			wantErr:      false,
		},
		{
			name: "negative 1 - invalid range, end before start",
			args: args{
				url: "",
				params: apiSpec.GetQueryRangeParams{
					Query:   `rate(some_metrics{env="prod"})`,
					Start:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
					End:     time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
					Step:    "5m",
					Timeout: "5s",
				},
			},
			wantRespCode: http.StatusBadRequest,
			wantErr:      false,
		},
		{
			name: "positive 1 - range = 1 hour, step = 5 minutes",
			args: args{
				url: "",
				params: apiSpec.GetQueryRangeParams{
					Query:   `rate(some_metrics{env="prod"})`,
					Start:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
					End:     time.Now().Format(time.RFC3339),
					Step:    "5m",
					Timeout: "5s",
				},
			},
			wantRespCode: http.StatusOK,
			wantErr:      false,
			wantCount:    13,
		},
		{
			name: "positive 2 - big range (1 year) should be cut to 24 hours",
			args: args{
				url: "",
				params: apiSpec.GetQueryRangeParams{
					Query:   `rate(some_metrics{env="prod"})`,
					Start:   time.Date(2020, time.June, 10, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
					End:     time.Date(2021, time.June, 10, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
					Step:    "6h",
					Timeout: "5s",
				},
			},
			wantRespCode: 200,
			wantErr:      false,
			wantCount:    5,
		},
	}

	for _, tt := range tests {
		ttCase := tt
		t.Run(ttCase.name, func(t *testing.T) {
			t.Parallel()

			req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, ttCase.args.url, nil)
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})} // nolint:exhaustivestruct
			ctx := echo.New().NewContext(req, rec)

			s := handlers{} // nolint:exhaustivestruct

			if err := s.GetQueryRange(ctx, ttCase.args.params); (err != nil) != ttCase.wantErr {
				t.Errorf("GetQueryRange() error = %v, wantErr %v", err, ttCase.wantErr)
			}

			if rec.Code != ttCase.wantRespCode {
				t.Errorf("GetQueryRange() code = %v, want %v", rec.Code, ttCase.wantRespCode)
			}

			if ttCase.wantCount > 0 {
				respBody, _ := io.ReadAll(rec.Body)
				var respObj apiSpec.QueryRangeResponse
				if err := json.Unmarshal(respBody, &respObj); err != nil {
					t.Error(err)
				}

				if len(respObj.Data.Result[0].Values) != ttCase.wantCount {
					t.Errorf("GetQueryRange() items len = %v, want %v", len(respObj.Data.Result[0].Values), ttCase.wantCount)
				}
			}
		})
	}
}
