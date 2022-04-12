package prometheus

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_fakeMetricsRandGenerator_Load(t *testing.T) {
	t.Parallel()

	t.Run("idempotent", func(t *testing.T) {
		t.Parallel()

		var (
			query = "hello world"
			start = time.Now().Add(5 * time.Second)
			stop  = time.Now().Add(10 * time.Minute)
			step  = time.Minute
		)

		g := fakeMetricsRandGenerator{}

		for i := 0; i < 10; i++ {
			got1 := g.Load(query, start, stop, step)
			got2 := g.Load(query, start, stop, step)

			if reflect.DeepEqual(got1, got2) {
				t.Errorf("Load() has no diff %v and %v", got1, got2)
			}
		}
	})

	t.Run("negative - wrong step", func(t *testing.T) {
		t.Parallel()

		t.Run("eq 0", func(t *testing.T) {
			t.Parallel()

			var (
				query = "hello wildfowl"
				start = time.Date(1990, time.June, 10, 8, 45, 0o0, 0, time.UTC)
				stop  = time.Date(1990, time.June, 10, 9, 0, 0o0, 0, time.UTC)
				step  = 0 * time.Minute
			)

			g := fakeMetricsRandGenerator{}

			assert.Panics(t, func() {
				_ = g.Load(query, start, stop, step)
			})
		})

		t.Run("less than 0", func(t *testing.T) {
			t.Parallel()

			var (
				query = "hello wildfowl"
				start = time.Date(1990, time.June, 10, 8, 45, 0o0, 0, time.UTC)
				stop  = time.Date(1990, time.June, 10, 9, 0, 0o0, 0, time.UTC)
				step  = -1 * time.Minute
			)

			g := fakeMetricsRandGenerator{}

			assert.Panics(t, func() {
				_ = g.Load(query, start, stop, step)
			})
		})
	})
}
