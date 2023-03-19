package prometheus

import (
	"fmt"
	"time"
)

var (
	ErrZeroStep            = fmt.Errorf("step should not be zero")
	ErrNegativeStep        = fmt.Errorf("step should not be negative")
	ErrParseDurationFailed = fmt.Errorf("can't parse duration")
)

// FakeMetricsGenerator describe methods for metrics generator.
type FakeMetricsGenerator interface {
	Load(query string, start, stop time.Time, step time.Duration) []seriesPoint
}
