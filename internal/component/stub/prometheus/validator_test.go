package prometheus //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	stubprometheus "github.com/goforbroke1006/onix/api/stub_prometheus"
)

func Test_validator_GetQueryRange(t *testing.T) {
	t.Parallel()

	target := validator{}

	t.Run("step is not golang-style duration", func(t *testing.T) {
		t.Parallel()

		err := target.GetQueryRange(stubprometheus.GetQueryRangeParams{
			Query:   "hello",
			Start:   "2006-01-02T15:04:05+07:00",
			End:     "2006-01-02T15:04:05+07:00",
			Step:    "1h30m - WILDFOWL",
			Timeout: "10s",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "can't parse step: 1h30m - WILDFOWL: can't parse duration", err.Error())
	})

	t.Run("timeout is not golang-style duration", func(t *testing.T) {
		t.Parallel()

		err := target.GetQueryRange(stubprometheus.GetQueryRangeParams{
			Query:   "hello",
			Start:   "2006-01-02T15:04:05+07:00",
			End:     "2006-01-02T15:04:05+07:00",
			Step:    "1h30m",
			Timeout: "hello10s",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "can't parse timeout: hello10s: can't parse duration", err.Error())
	})
}
