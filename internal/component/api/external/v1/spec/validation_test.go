package spec

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReportCompareRequest_Validate(t *testing.T) { // nolint:funlen
	t.Parallel()

	type fields struct {
		Service   string
		Source    string
		TagOne    string
		TagTwo    string
		TimeRange ReportTimeRange
	}
	tests := []struct {
		name   string
		fields fields
		want   []error
	}{
		{
			name: "negative - all fields valid",
			fields: fields{
				Service:   "acme.get-user-ingo",
				Source:    "prometheus.prod.env",
				TagOne:    "v1.0.0",
				TagTwo:    "v2.0.0",
				TimeRange: "1d",
			},
			want: nil,
		},
		{
			name: "positive 1 - service empty",
			fields: fields{
				Service:   "",
				Source:    "prometheus.prod.env",
				TagOne:    "v1.0.0",
				TagTwo:    "v2.0.0",
				TimeRange: "1d",
			},
			want: []error{
				errors.New("`service` should not be empty"),
			},
		},
		{
			name: "positive 2 - all fields are empty",
			fields: fields{
				Service:   "",
				Source:    "",
				TagOne:    "",
				TagTwo:    "",
				TimeRange: "",
			},
			want: []error{
				errors.New("`service` should not be empty"),
				errors.New("`source` should not be empty"),
				errors.New("`tag_one` should not be empty"),
				errors.New("`tag_two` should not be empty"),
				errors.New("`time_range` should not be empty"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := ReportCompareRequest{
				Service:   tt.fields.Service,
				Source:    tt.fields.Source,
				TagOne:    tt.fields.TagOne,
				TagTwo:    tt.fields.TagTwo,
				TimeRange: tt.fields.TimeRange,
			}

			got := r.Validate()

			assert.Equal(t, len(got), len(tt.want))
			for idx := range tt.want {
				assert.Equal(t, tt.want[idx].Error(), got[idx].Error())
			}
		})
	}
}
