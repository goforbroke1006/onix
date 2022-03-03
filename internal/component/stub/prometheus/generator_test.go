package prometheus

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_fakeMetricsRandGenerator_Load(t *testing.T) {
	type args struct {
		query string
		start time.Time
		stop  time.Time
		step  time.Duration
	}

	g := fakeMetricsRandGenerator{}
	a := args{query: "hello world", start: time.Time{}, stop: time.Time{}, step: time.Minute}

	for i := 0; i < 10; i++ {
		got1 := g.Load(a.query, a.start, a.stop, a.step)
		got2 := g.Load(a.query, a.start, a.stop, a.step)
		if reflect.DeepEqual(got1, got2) {
			t.Errorf("Load() has no diff %v and %v", got1, got2)
		}
	}
}

func Test_fakeMetricsIdempotentGenerator_Load(t *testing.T) {
	type args struct {
		query string
		start time.Time
		stop  time.Time
		step  time.Duration
	}
	tests := []struct {
		name string
		args args
		want []seriesPoint
	}{
		{
			name: "positive - closed range = 1 item in result",
			args: args{query: "hello world", start: time.Time{}, stop: time.Time{}, step: time.Minute},
			want: []seriesPoint{{-62135596800, 0.5007022298180581}},
		},
		{
			name: "negative - invalid range = no result",
			args: args{query: "hello kitty", start: time.Now().Add(time.Minute), stop: time.Now().Add(-1 * time.Second), step: 5 * time.Second},
			want: nil,
		},
		{
			name: "positive - some range 15 min",
			args: args{
				query: "hello wildfowl",
				start: time.Date(1990, time.June, 10, 8, 45, 00, 0, time.UTC),
				stop:  time.Date(1990, time.June, 10, 9, 0, 00, 0, time.UTC),
				step:  5 * time.Minute,
			},
			want: []seriesPoint{
				{645007500, 0.5452772128272665},
				{645007800, 0.9309666611451856},
				{645008100, 0.32090824169586474},
				{645008400, 0.6830072977477972},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := fakeMetricsIdempotentGenerator{}
			if got := g.Load(tt.args.query, tt.args.start, tt.args.stop, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fakeMetricsIdempotentGenerator_hash(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "positive - empty string, got default seed",
			args: args{query: ""},
			want: 1579219711395102,
		},
		{
			name: "positive - str 1",
			args: args{query: "hello world"},
			want: 6144619784920104,
		},
		{
			name: "positive 2 - more real sample",
			args: args{query: `histogram_quantile(0.95, sum(increase(api_request_count{environment="prod",instrument="one"}[15m])) by (le))`},
			want: 3845172339482296,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := fakeMetricsIdempotentGenerator{}

			got := g.hash(tt.args.query)

			if got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}

			gotToStr := fmt.Sprintf("%d", got)
			if len(gotToStr) != 16 {
				t.Errorf("hash() length = %v, want %v", len(gotToStr), 16)
			}
		})
	}
}
