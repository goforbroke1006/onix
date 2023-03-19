package spec

import (
	"time"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/internal/component/stub/prometheus"
)

func (p GetQueryRangeParams) Validate() []error {
	var errs []error

	if len(p.Query) == 0 {
		errs = append(errs, errors.New("`query` should not be empty"))
	}

	var (
		start time.Time
		stop  time.Time
		err   error
	)

	if start, err = prometheus.CanParseTime(p.Start); err != nil {
		errs = append(errs, errors.New("`start` contains invalid date-time"))
	}

	if stop, err = prometheus.CanParseTime(p.End); err != nil {
		errs = append(errs, errors.New("`end` contains invalid date-time"))
	}

	if _, err = prometheus.CanParseDuration(p.Step); err != nil {
		errs = append(errs, errors.New("`step` contains invalid date-time"))
	}

	if _, err = prometheus.CanParseDuration(string(p.Timeout)); err != nil {
		errs = append(errs, errors.New("`timeout` contains invalid date-time"))
	}

	if stop.Before(start) {
		errs = append(errs, errors.New("wrong time range"))
	}

	return errs
}
