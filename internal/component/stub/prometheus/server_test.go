package prometheus

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func Test_server_GetHealthz(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name         string
		args         args
		wantRespCode int
		wantErr      bool
	}{
		{
			name:         "ok for any request",
			args:         args{},
			wantRespCode: 200,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.args.url, nil)
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
			ctx := echo.New().NewContext(req, rec)

			s := server{}
			err := s.GetHealthz(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetHealthz() error = %v, wantErr %v", err, tt.wantErr)
			}
			if rec.Code != tt.wantRespCode {
				t.Errorf("GetHealthz() status code, got = %v, want %v", rec.Code, tt.wantRespCode)
			}
		})
	}
}

func Test_server_GetQueryRange(t *testing.T) {
	type args struct {
		url    string
		params GetQueryRangeParams
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
				params: GetQueryRangeParams{
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
				params: GetQueryRangeParams{
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
				params: GetQueryRangeParams{
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
				params: GetQueryRangeParams{
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
				params: GetQueryRangeParams{
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
				params: GetQueryRangeParams{
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
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.args.url, nil)
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
			ctx := echo.New().NewContext(req, rec)

			s := server{}
			if err := s.GetQueryRange(ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("GetQueryRange() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantRespCode {
				t.Errorf("GetQueryRange() code = %v, want %v", rec.Code, tt.wantRespCode)
			}

			if tt.wantCount > 0 {
				respBody, _ := ioutil.ReadAll(rec.Body)
				respObj := QueryRangeResponse{}
				if err := json.Unmarshal(respBody, &respObj); err != nil {
					t.Error(err)
				}

				if len(respObj.Data.Result[0].Values) != tt.wantCount {
					t.Errorf("GetQueryRange() items len = %v, want %v", len(respObj.Data.Result[0].Values), tt.wantCount)
				}
			}
		})
	}
}

func Test_server_canParseTime(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "negative - empty",
			args:    args{str: ""},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "negative - date only",
			args:    args{str: "2022-02-01"},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "positive - RFC 3339",
			args:    args{str: "2022-02-01T12:34:56Z"},
			want:    time.Date(2022, time.February, 1, 12, 34, 56, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "positive - unix timestamp",
			args:    args{str: "1643718896"},
			want:    time.Date(2022, time.February, 1, 12, 34, 56, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server{}
			got, err := s.canParseTime(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("canParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("canParseTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_canParseDuration(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "negative - invalid",
			args:    args{str: "hello world"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "positive - float",
			args:    args{str: "3.14"},
			want:    3140 * time.Millisecond,
			wantErr: false,
		},
		{
			name:    "positive - duration in go style",
			args:    args{str: "5m"},
			want:    5 * time.Minute,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server{}
			got, err := s.canParseDuration(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("canParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("canParseDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
