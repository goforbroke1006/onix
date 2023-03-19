package prometheus //nolint:testpackage

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/domain"
)

func Test_fakeMetricsIdempotentGenerator_Load(t *testing.T) { //nolint:funlen
	t.Parallel()

	type args struct {
		query string
		start time.Time
		stop  time.Time
		step  time.Duration
	}

	type testCase struct {
		name string
		args args
		want []domain.SeriesPoint
	}

	tests := []testCase{
		{
			name: "positive - closed range = 1 item in result",
			args: args{
				query: "hello world",
				start: time.Time{},
				stop:  time.Time{},
				step:  time.Minute,
			},
			want: []domain.SeriesPoint{{-62135596800, 0.5007022298180581}},
		},
		{
			name: "negative - invalid range = no result",
			args: args{
				query: "hello kitty",
				start: time.Now().Add(time.Minute),
				stop:  time.Now().Add(-1 * time.Second),
				step:  5 * time.Second,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		ttCase := tt
		func(ttCase testCase) {
			t.Run(ttCase.name, func(t *testing.T) {
				t.Parallel()

				g := fakeMetricsIdempotentGenerator{}
				got := g.Load(ttCase.args.query, ttCase.args.start, ttCase.args.stop, ttCase.args.step)
				if !reflect.DeepEqual(got, ttCase.want) {
					t.Errorf("Load() = %v, want %v", got, ttCase.want)
				}
			})
		}(ttCase)
	}

	t.Run("no regress", func(t *testing.T) {
		t.Parallel()

		var (
			query = "hello wildfowl"
			start = time.Date(1990, time.June, 10, 8, 45, 0o0, 0, time.UTC)
			stop  = time.Date(1990, time.June, 10, 9, 0, 0o0, 0, time.UTC)
			step  = 5 * time.Minute
		)
		expected := []domain.SeriesPoint{
			{645007500, 0.5452772128272665},
			{645007800, 0.9309666611451856},
			{645008100, 0.32090824169586474},
			{645008400, 0.6830072977477972},
		}

		g := fakeMetricsIdempotentGenerator{}
		got := g.Load(query, start, stop, step)

		assert.Equal(t, expected, got)
	})

	t.Run("negative - wrong step", func(t *testing.T) {
		t.Parallel()

		var (
			query = "hello wildfowl"
			start = time.Date(1990, time.June, 10, 8, 45, 0o0, 0, time.UTC)
			stop  = time.Date(1990, time.June, 10, 9, 0, 0o0, 0, time.UTC)
		)

		generator := fakeMetricsIdempotentGenerator{}

		t.Run("eq 0", func(t *testing.T) {
			t.Parallel()

			assert.Panics(t, func() {
				_ = generator.Load(query, start, stop, 0)
			})
		})

		t.Run("less than 0", func(t *testing.T) {
			t.Parallel()

			assert.Panics(t, func() {
				_ = generator.Load(query, start, stop, -1*time.Minute)
			})
		})
	})
}

func Test_fakeMetricsIdempotentGenerator_hash(t *testing.T) { //nolint:tparallel
	t.Parallel()

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
			args: args{
				query: `
histogram_quantile(
  0.95, 
  sum(
    increase(api_request_count{
      environment="prod",instrument="one"
    }[15m])
  ) by (le)
)`,
			},
			want: 5217508800143032,
		},
	}

	for _, tt := range tests { //nolint:paralleltest
		ttCase := tt
		t.Run(ttCase.name, func(t *testing.T) {
			g := fakeMetricsIdempotentGenerator{}

			got := g.hash(ttCase.args.query)

			if got != ttCase.want {
				t.Errorf("hash() = %v, want %v", got, ttCase.want)
			}

			gotToStr := fmt.Sprintf("%d", got)
			if len(gotToStr) != 16 {
				t.Errorf("hash() length = %v, want %v", len(gotToStr), 16)
			}
		})
	}
}
