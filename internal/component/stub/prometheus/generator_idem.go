package prometheus

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type seriesPoint struct {
	timestamp int64
	value     float64
}

// FakeMetricsGenerator describe methods for metrics generator.
type FakeMetricsGenerator interface {
	Load(query string, start, stop time.Time, step time.Duration) []seriesPoint
}

var _ FakeMetricsGenerator = &fakeMetricsIdempotentGenerator{}

type fakeMetricsIdempotentGenerator struct{}

func (g fakeMetricsIdempotentGenerator) Load(query string, start, stop time.Time, step time.Duration) []seriesPoint {
	hash := g.hash(query)

	if step == 0 {
		panic(ErrZeroStep)
	}

	if step < 0 {
		panic(ErrNegativeStep)
	}

	var result []seriesPoint

	current := start

	for current.Before(stop) || current.Equal(stop) {
		rg := rand.New(rand.NewSource(hash * current.UnixNano()))
		f := rg.Float64() // nolint:gosec

		result = append(result, seriesPoint{
			timestamp: current.Unix(),
			value:     f,
		})
		current = current.Add(step)
	}

	return result
}

const (
	defaultIdempotentSeed  = 123
	defaultIdempotentBoost = 12
)

// hash generates int64 16-digit number for provided query.
func (g fakeMetricsIdempotentGenerator) hash(query string) int64 {
	result := int64(defaultIdempotentSeed)

	const expectedLen = 16

	// mix up with string content
	for index, letter := range query {
		result += int64(letter) * int64(index)
	}

	// raise digits count in result
	const tenBase = 10
	bound := int64(math.Pow(tenBase, expectedLen))

	for result < bound {
		result *= defaultIdempotentBoost
	}

	result *= defaultIdempotentBoost

	// cut extra digits
	str := fmt.Sprintf("%d", result)[:expectedLen]

	const (
		base    = 10
		bitSize = 64
	)

	result, _ = strconv.ParseInt(str, base, bitSize)

	return result
}
