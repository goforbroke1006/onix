package prometheus

import (
	"math/rand"
	"time"

	"github.com/goforbroke1006/onix/domain"
)

func NewFakeMetricsRandGenerator() domain.FakeMetricsGenerator {
	return &fakeMetricsRandGenerator{
		randomizer: rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}
}

var _ domain.FakeMetricsGenerator = (*fakeMetricsRandGenerator)(nil)

type fakeMetricsRandGenerator struct {
	randomizer *rand.Rand
}

// Load returns totally random data for each invocation.
func (g fakeMetricsRandGenerator) Load(query string, start, stop time.Time, step time.Duration) []domain.SeriesPoint {
	_ = query

	if step == 0 {
		panic(domain.ErrZeroStep)
	}

	if step < 0 {
		panic(domain.ErrNegativeStep)
	}

	result := make([]domain.SeriesPoint, 0, stop.Sub(start)/step+1)
	current := start

	for current.Before(stop) || current.Equal(stop) {
		f := g.randomizer.Float64()

		result = append(result, domain.SeriesPoint{
			Timestamp: current.Unix(),
			Value:     f,
		})
		current = current.Add(step)
	}

	return result
}
