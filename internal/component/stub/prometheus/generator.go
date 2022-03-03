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

type FakeMetricsGenerator interface {
	Load(query string, start, stop time.Time, step time.Duration) []seriesPoint
}

var (
	_ FakeMetricsGenerator = &fakeMetricsRandGenerator{}
)

type fakeMetricsRandGenerator struct {
}

// Load returns totally random data for each invocation
func (g fakeMetricsRandGenerator) Load(query string, start, stop time.Time, step time.Duration) []seriesPoint {
	rand.Seed(time.Now().UnixNano())

	if step == 0 {
		panic(fmt.Errorf("step should not be zero"))
	}
	if step < 0 {
		panic(fmt.Errorf("step should not be negative"))
	}

	result := make([]seriesPoint, 0, stop.Sub(start)/step+1)
	current := start
	for current.Before(stop) || current.Equal(stop) {
		result = append(result, seriesPoint{
			timestamp: current.Unix(),
			value:     rand.Float64(),
		})
		current = current.Add(step)
	}
	return result
}

var (
	_ FakeMetricsGenerator = &fakeMetricsIdempotentGenerator{}
)

type fakeMetricsIdempotentGenerator struct {
}

func (g fakeMetricsIdempotentGenerator) Load(query string, start, stop time.Time, step time.Duration) []seriesPoint {
	rand.Seed(g.hash(query))

	if step == 0 {
		panic(fmt.Errorf("step should not be zero"))
	}
	if step < 0 {
		panic(fmt.Errorf("step should not be negative"))
	}

	var result []seriesPoint
	current := start
	for current.Before(stop) || current.Equal(stop) {
		result = append(result, seriesPoint{
			timestamp: current.Unix(),
			value:     rand.Float64(),
		})
		current = current.Add(step)
	}
	return result
}

const defaultIdempotentSeed = 123
const defaultIdempotentBoost = 12

// hash generates int64 16-digit number for provided query
func (g fakeMetricsIdempotentGenerator) hash(query string) int64 {
	result := int64(defaultIdempotentSeed)
	const expectedLen = 16

	// mix up with string content
	for index, letter := range query {
		result += int64(letter) * int64(index)
	}

	// raise digits count in result
	bound := int64(math.Pow(10, expectedLen))
	for result < bound {
		//next := result * result
		//if next == result {
		//	break
		//}
		//result = next
		result *= defaultIdempotentBoost
	}

	result *= defaultIdempotentBoost

	// cut extra digits
	str := fmt.Sprintf("%d", result)[:expectedLen]
	result, _ = strconv.ParseInt(str, 10, 64)

	return result
}
