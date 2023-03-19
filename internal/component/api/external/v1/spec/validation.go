package spec

import (
	"github.com/pkg/errors"
)

func (r ReportCompareRequest) Validate() []error {
	var errs []error

	if r.Service == "" {
		errs = append(errs, errors.New("`service` should not be empty"))
	}
	if r.Source == "" {
		errs = append(errs, errors.New("`source` should not be empty"))
	}
	if r.TagOne == "" {
		errs = append(errs, errors.New("`tag_one` should not be empty"))
	}
	if r.TagTwo == "" {
		errs = append(errs, errors.New("`tag_two` should not be empty"))
	}
	if r.TimeRange == "" {
		errs = append(errs, errors.New("`time_range` should not be empty"))
	}

	return errs
}
