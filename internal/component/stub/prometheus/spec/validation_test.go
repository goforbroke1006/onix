package spec //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validator_GetQueryRange(t *testing.T) {
	t.Parallel()

	t.Run("step is not golang-style duration", func(t *testing.T) {
		t.Parallel()

		errs := GetQueryRangeParams{
			Query:   "hello",
			Start:   "2006-01-02T15:04:05+07:00",
			End:     "2006-01-02T15:04:05+07:00",
			Step:    "1h30m - WILDFOWL",
			Timeout: "10s",
		}.Validate()

		assert.NotNil(t, errs)
		assert.Len(t, errs, 1)
		assert.Equal(t, "`step` contains invalid date-time", errs[0].Error())
	})

	t.Run("timeout is not golang-style duration", func(t *testing.T) {
		t.Parallel()

		errs := GetQueryRangeParams{
			Query:   "hello",
			Start:   "2006-01-02T15:04:05+07:00",
			End:     "2006-01-02T15:04:05+07:00",
			Step:    "1h30m",
			Timeout: "hello10s",
		}.Validate()
		assert.NotNil(t, errs)
		assert.Len(t, errs, 1)
		assert.Equal(t, "`timeout` contains invalid date-time", errs[0].Error())
	})
}
