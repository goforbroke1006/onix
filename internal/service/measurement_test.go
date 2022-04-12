package service // nolint:testpackage

import (
	"context"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/domain"
	mockRepo "github.com/goforbroke1006/onix/internal/repository/mocks"
	mocksMP "github.com/goforbroke1006/onix/internal/service/metricsprovider/mocks"
)

func TestNewMeasurementService(t *testing.T) {
	t.Parallel()

	actual := NewMeasurementService(nil)
	assert.NotNil(t, actual)
}

func Test_measurementService_getTimePoints(t *testing.T) { // nolint:funlen
	t.Parallel()

	type fields struct{}

	type args struct {
		from time.Time
		till time.Time
		step time.Duration
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []time.Time
	}{
		{
			name:   "zero interval",
			fields: fields{},
			args: args{
				from: time.Time{},
				till: time.Time{},
				step: 15 * time.Minute,
			},
			want: []time.Time{
				{},
			},
		},
		{
			name:   "usual interval",
			fields: fields{},
			args: args{
				from: time.Date(2022, time.March, 8, 13, 0, 0, 0, time.UTC),
				till: time.Date(2022, time.March, 8, 14, 0, 0, 0, time.UTC),
				step: 15 * time.Minute,
			},
			want: []time.Time{
				time.Date(2022, time.March, 8, 13, 0, 0, 0, time.UTC),
				time.Date(2022, time.March, 8, 13, 15, 0, 0, time.UTC),
				time.Date(2022, time.March, 8, 13, 30, 0, 0, time.UTC),
				time.Date(2022, time.March, 8, 13, 45, 0, 0, time.UTC),
				time.Date(2022, time.March, 8, 14, 0o0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			m := measurementService{} // nolint:exhaustivestruct
			got := m.getTimePoints(testCase.args.from, testCase.args.till, testCase.args.step)
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("getTimePoints() = %v, want %v", got, testCase.want)
			}
		})
	}

	t.Run("debug", func(t *testing.T) {
		var (
			from = time.Unix(1642877700, 0) // Sat Jan 22 2022 18:55:00 GMT+0000
			till = time.Unix(1642895700, 0) // Sat Jan 22 2022 23:55:00 GMT+0000
		)

		var m measurementService
		got := m.getTimePoints(from, till, 5*time.Minute)
		assert.Equal(t, 12*5+1, len(got))
	})
}

func TestMeasurementService_GetOrPull(t *testing.T) { // nolint:funlen,maintidx
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("on DB failure", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		measurementRepository := mockRepo.NewMockMeasurementRepository(mockCtrl)
		svc := measurementService{
			measurementRepo:  measurementRepository,
			createProviderFn: nil,
		}

		measurementRepository.EXPECT().GetForPoints(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nil,
			errors.New("some fake error"),
		)

		source := domain.Source{ // nolint:exhaustivestruct
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{} // nolint:exhaustivestruct

		from := time.Date(2022, time.March, 8, 13, 0, 0, 0, time.UTC)
		till := time.Date(2022, time.March, 8, 14, 0, 0, 0, time.UTC)

		_, err := svc.GetOrPull(context.TODO(), source, criteria, from, till, 15*time.Minute)
		assert.NotNil(t, err)
	})

	t.Run("just return from DB if found series", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		measurementRepository := mockRepo.NewMockMeasurementRepository(mockCtrl)
		svc := measurementService{
			measurementRepo:  measurementRepository,
			createProviderFn: nil,
		}

		measurementRepository.EXPECT().GetForPoints(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			[]domain.MeasurementRow{
				{Moment: time.Time{}, Value: 1},
				{Moment: time.Time{}, Value: 2},
				{Moment: time.Time{}, Value: 5},
				{Moment: time.Time{}, Value: 5},
				{Moment: time.Time{}, Value: 5},
			},
			nil,
		)

		source := domain.Source{ // nolint:exhaustivestruct
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{} // nolint:exhaustivestruct

		from := time.Date(2022, time.March, 8, 13, 0, 0, 0, time.UTC)
		till := time.Date(2022, time.March, 8, 14, 0, 0, 0, time.UTC)

		actual, err := svc.GetOrPull(context.TODO(), source, criteria, from, till, 15*time.Minute)
		if err != nil {
			t.Fatal(err)
		}

		expected := []domain.MeasurementRow{
			{Moment: time.Time{}, Value: 1},
			{Moment: time.Time{}, Value: 2},
			{Moment: time.Time{}, Value: 5},
			{Moment: time.Time{}, Value: 5},
			{Moment: time.Time{}, Value: 5},
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("get series error", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		measurementRepository := mockRepo.NewMockMeasurementRepository(mockCtrl)
		svc := measurementService{
			measurementRepo: measurementRepository,
			createProviderFn: func(source domain.Source) domain.MetricsProvider {
				provider := mocksMP.NewMockMetricsProvider(mockCtrl)
				provider.
					EXPECT().
					LoadSeries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(
						nil,
						errors.New("some fake errors"),
					)

				return provider
			},
		}

		measurementRepository.EXPECT().GetForPoints(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			[]domain.MeasurementRow{}, // measurement table empty
			nil,
		)

		source := domain.Source{ // nolint:exhaustivestruct
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{} // nolint:exhaustivestruct

		from := time.Date(2022, time.March, 8, 13, 0, 0, 0, time.UTC)
		till := time.Date(2022, time.March, 8, 14, 0, 0, 0, time.UTC)

		_, err := svc.GetOrPull(context.TODO(), source, criteria, from, till, 15*time.Minute)
		assert.NotNil(t, err)
	})

	t.Run("get series from provider but can't store in DB", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		measurementRepository := mockRepo.NewMockMeasurementRepository(mockCtrl)
		svc := measurementService{
			measurementRepo: measurementRepository,
			createProviderFn: func(source domain.Source) domain.MetricsProvider {
				provider := mocksMP.NewMockMetricsProvider(mockCtrl)
				provider.
					EXPECT().
					LoadSeries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(
						[]domain.SeriesItem{
							{Timestamp: time.Time{}, Value: 0},
							{Timestamp: time.Time{}, Value: 0},
							{Timestamp: time.Time{}, Value: 0},
						},
						nil,
					)

				return provider
			},
		}

		measurementRepository.EXPECT().GetForPoints(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			[]domain.MeasurementRow{}, // measurement table empty
			nil,
		)
		measurementRepository.EXPECT().StoreBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			errors.New("some fake error"),
		)

		source := domain.Source{ // nolint:exhaustivestruct
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{} // nolint:exhaustivestruct

		from := time.Date(2022, time.March, 8, 13, 0, 0, 0, time.UTC)
		till := time.Date(2022, time.March, 8, 14, 0, 0, 0, time.UTC)

		_, err := svc.GetOrPull(context.TODO(), source, criteria, from, till, 15*time.Minute)
		assert.NotNil(t, err)
	})

	t.Run("get series from provider and store in DB", func(t *testing.T) {
		t.Parallel()

		measurementRepository := mockRepo.NewMockMeasurementRepository(mockCtrl)
		measurementRepository.EXPECT().GetForPoints(gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]domain.MeasurementRow{}, nil) // measurement table empty
		measurementRepository.EXPECT().StoreBatch(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil) // store batch OK

		svc := measurementService{
			measurementRepo: measurementRepository,
			createProviderFn: func(source domain.Source) domain.MetricsProvider {
				provider := mocksMP.NewMockMetricsProvider(mockCtrl)
				provider.
					EXPECT().
					LoadSeries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(
						[]domain.SeriesItem{
							{Timestamp: parseDateTimeToUTC("2022-03-07 13:00:00"), Value: 1},
							{Timestamp: parseDateTimeToUTC("2022-03-07 13:15:00"), Value: 2},
							{Timestamp: parseDateTimeToUTC("2022-03-07 14:00:00"), Value: 5},
						},
						nil,
					)

				return provider
			},
		}

		source := domain.Source{ // nolint:exhaustivestruct
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{} // nolint:exhaustivestruct

		from := time.Date(2022, time.March, 7, 13, 0, 0, 0, time.UTC)
		till := time.Date(2022, time.March, 7, 14, 0, 0, 0, time.UTC)

		actual, err := svc.GetOrPull(context.TODO(), source, criteria, from, till, 15*time.Minute)
		if err != nil {
			t.Fatal(err)
		}

		expected := []domain.MeasurementRow{
			{Moment: parseDateTimeToUTC("2022-03-07 13:00:00"), Value: 1},
			{Moment: parseDateTimeToUTC("2022-03-07 13:15:00"), Value: 2},
			{Moment: parseDateTimeToUTC("2022-03-07 13:30:00"), Value: 0},
			{Moment: parseDateTimeToUTC("2022-03-07 13:45:00"), Value: 0},
			{Moment: parseDateTimeToUTC("2022-03-07 14:00:00"), Value: 5},
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("replace missed points with zeros", func(t *testing.T) {
		t.Parallel()

		measurementRepo := mockRepo.NewMockMeasurementRepository(mockCtrl)
		measurementRepo.EXPECT().GetForPoints(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		measurementRepo.EXPECT().StoreBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		provider := mocksMP.NewMockMetricsProvider(mockCtrl)
		provider.EXPECT().LoadSeries(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]domain.SeriesItem{
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:00:00"), Value: 10},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:05:00"), Value: 11},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:10:00"), Value: 12},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:15:00"), Value: 13},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:20:00"), Value: 14},
				// {Timestamp: parseDateTimeToUTC("2022-01-22 10:25:00"), Value: 15},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:30:00"), Value: 16},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:35:00"), Value: 17},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:40:00"), Value: 18},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:45:00"), Value: 19},
				// {Timestamp: parseDateTimeToUTC("2022-01-22 10:50:00"), Value: 20},
				{Timestamp: parseDateTimeToUTC("2022-01-22 10:55:00"), Value: 21},
				{Timestamp: parseDateTimeToUTC("2022-01-22 11:00:00"), Value: 22},
			}, nil)

		svc := measurementService{
			measurementRepo: measurementRepo,
			createProviderFn: func(_ domain.Source) domain.MetricsProvider {
				return provider
			},
		}
		points, err := svc.GetOrPull(
			context.TODO(),
			domain.Source{ID: 0, Title: "", Type: "", Address: ""},
			domain.Criteria{ID: 0, Service: "", Title: "", Selector: "", ExpectedDir: "", GroupingInterval: 0},
			parseDateTimeToUTC("2022-01-22 10:00:00"),
			parseDateTimeToUTC("2022-01-22 11:00:00"),
			5*time.Minute,
		)
		assert.Nil(t, err)
		assert.Len(t, points, 13)

		assert.Equal(t, 10.0, points[0].Value)
		assert.Equal(t, 11.0, points[1].Value)
		assert.Equal(t, 0.0, points[5].Value)
		assert.Equal(t, 0.0, points[10].Value)
		assert.Equal(t, 21.0, points[11].Value)
		assert.Equal(t, 22.0, points[12].Value)
	})
}

func parseDateTimeToUTC(str string) time.Time {
	re := regexp.MustCompile(`([\d]{4})-([\d]{2})-([\d]{2}) ([\d]{2}):([\d]{2}):([\d]{2})`)
	submatch := re.FindStringSubmatch(str)

	var (
		year, _  = strconv.ParseInt(submatch[1], 10, 64)
		month, _ = strconv.ParseInt(submatch[2], 10, 64)
		day, _   = strconv.ParseInt(submatch[3], 10, 64)

		hours, _   = strconv.ParseInt(submatch[4], 10, 64)
		minutes, _ = strconv.ParseInt(submatch[5], 10, 64)
		seconds, _ = strconv.ParseInt(submatch[6], 10, 64)
	)

	return time.Date(
		int(year), time.Month(month), int(day),
		int(hours), int(minutes), int(seconds), 0,
		time.UTC)
}
