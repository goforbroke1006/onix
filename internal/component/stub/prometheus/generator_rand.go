package prometheus

import (
	"math/rand"
	"time"
)

var _ FakeMetricsGenerator = &fakeMetricsRandGenerator{}

type fakeMetricsRandGenerator struct{}

// Load returns totally random data for each invocation.
func (g fakeMetricsRandGenerator) Load(query string, start, stop time.Time, step time.Duration) []seriesPoint {
	rand.Seed(time.Now().UnixNano())

	if step == 0 {
		panic(ErrZeroStep)
	}

	if step < 0 {
		panic(ErrNegativeStep)
	}

	result := make([]seriesPoint, 0, stop.Sub(start)/step+1)
	current := start

	for current.Before(stop) || current.Equal(stop) {
		f := rand.Float64() // nolint:gosec

		result = append(result, seriesPoint{
			timestamp: current.Unix(),
			value:     f,
		})
		current = current.Add(step)
	}

	return result
}
