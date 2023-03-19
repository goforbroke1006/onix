package prometheus

import (
	"math/rand"
	"time"
)

func NewFakeMetricsRandGenerator() FakeMetricsGenerator {
	return &fakeMetricsRandGenerator{
		randomizer: rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}
}

var _ FakeMetricsGenerator = (*fakeMetricsRandGenerator)(nil)

type fakeMetricsRandGenerator struct {
	randomizer *rand.Rand
}

// Load returns totally random data for each invocation.
func (g fakeMetricsRandGenerator) Load(query string, start, stop time.Time, step time.Duration) []seriesPoint {
	_ = query

	if step == 0 {
		panic(ErrZeroStep)
	}

	if step < 0 {
		panic(ErrNegativeStep)
	}

	result := make([]seriesPoint, 0, stop.Sub(start)/step+1)
	current := start

	for current.Before(stop) || current.Equal(stop) {
		f := g.randomizer.Float64()

		result = append(result, seriesPoint{
			timestamp: current.Unix(),
			value:     f,
		})
		current = current.Add(step)
	}

	return result
}
