package prometheus // nolint:testpackage

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_server_canParseTime(t *testing.T) { // nolint:funlen
	t.Parallel()

	type args struct {
		str string
	}

	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "negative 1 - empty string",
			args:    args{str: ""},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "negative 2 - invalid string",
			args:    args{str: "hello world"},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "negative 3 - date only",
			args:    args{str: "2022-02-01"},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "positive 1 - int",
			args:    args{str: "1647270683"},
			want:    time.Date(2022, time.March, 14, 15, 11, 23, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "positive 2 - rfc 3339 utc",
			args:    args{str: "2022-03-14T15:11:23Z"},
			want:    time.Date(2022, time.March, 14, 15, 11, 23, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "positive 3 - rfc 3339 another timezone",
			args:    args{str: "2022-03-14T18:11:23.0+03:00"},
			want:    time.Date(2022, time.March, 14, 15, 11, 23, 0, time.UTC),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		ttCase := tt
		t.Run(ttCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := canParseTime(ttCase.args.str)
			if (err != nil) != ttCase.wantErr {
				t.Errorf("canParseTime() error = %v, wantErr %v", err, ttCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, ttCase.want) {
				t.Errorf("canParseTime() got = %v, want %v", got, ttCase.want)
			}
		})
	}
}

func Test_server_canParseTime_withBrokenDigitRegex(t *testing.T) {
	t.Parallel()

	onlyNumbersRegex = regexp.MustCompile(`^[\w]+$`)

	_, err := canParseTime("123hello")
	assert.NotNil(t, err)
	assert.Equal(t, "can't parse integer: strconv.ParseInt: parsing \"123hello\": invalid syntax", err.Error())
}

func Test_server_canParseDuration(t *testing.T) {
	t.Parallel()

	type args struct {
		str string
	}

	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "negative - invalid",
			args:    args{str: "hello world"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "positive - float",
			args:    args{str: "3.14"},
			want:    3140 * time.Millisecond,
			wantErr: false,
		},
		{
			name:    "positive - duration in go style",
			args:    args{str: "5m"},
			want:    5 * time.Minute,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		ttCase := tt
		t.Run(ttCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := canParseDuration(ttCase.args.str)
			if (err != nil) != ttCase.wantErr {
				t.Errorf("canParseDuration() error = %v, wantErr %v", err, ttCase.wantErr)

				return
			}
			if got != ttCase.want {
				t.Errorf("canParseDuration() got = %v, want %v", got, ttCase.want)
			}
		})
	}
}
