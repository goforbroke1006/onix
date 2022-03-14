package prometheus

import (
	"time"

	"github.com/pkg/errors"

	apiSpec "github.com/goforbroke1006/onix/api/stub_prometheus"
)

type validator struct{}

func (v validator) GetQueryRange(params apiSpec.GetQueryRangeParams) error {
	if len(params.Query) == 0 {
		return errors.New("empty")
	}

	var (
		start time.Time
		stop  time.Time
		err   error
	)

	if start, err = canParseTime(params.Start); err != nil {
		return err
	}

	if stop, err = canParseTime(params.End); err != nil {
		return err
	}

	if _, err = canParseDuration(params.Step); err != nil {
		return err
	}

	if _, err = canParseDuration(string(params.Timeout)); err != nil {
		return err
	}

	if stop.Before(start) {
		return errors.New("wrong time range")
	}

	return nil
}
