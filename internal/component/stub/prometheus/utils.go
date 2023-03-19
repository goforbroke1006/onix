package prometheus

import (
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/goforbroke1006/onix/domain"
)

var onlyNumbersRegex = regexp.MustCompile(`^[\d]+$`)

func CanParseTime(str string) (time.Time, error) {
	if onlyNumbersRegex.Match([]byte(str)) {
		const (
			base    = 10
			bitSize = 64
		)

		var (
			unix int64
			err  error
		)

		if unix, err = strconv.ParseInt(str, base, bitSize); err != nil {
			return time.Time{}, errors.Wrap(err, "can't parse integer")
		}

		res := time.Unix(unix, 0).UTC()

		return res, nil
	}

	res, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "can't parse by RFC 3339")
	}

	return res.UTC(), nil
}

func CanParseDuration(str string) (time.Duration, error) {
	const bitsSize = 64
	float, err := strconv.ParseFloat(str, bitsSize)

	if err == nil {
		const (
			tenBase = 10
			degree  = 9
		)

		toNanoSeconds := math.Pow(tenBase, degree)

		return time.Duration(int(float*toNanoSeconds)) * time.Nanosecond, nil
	}

	duration, err := time.ParseDuration(str)
	if err == nil {
		return duration, nil
	}

	return 0, errors.Wrap(domain.ErrParseDurationFailed, str)
}
