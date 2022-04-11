package service // nolint:testpackage

import (
	"context"
	"reflect"
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

func TestMeasurementService_GetOrPull(t *testing.T) { // nolint:funlen
	t.Parallel()

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
							{Timestamp: time.Time{}, Value: 1},
							{Timestamp: time.Time{}, Value: 2},
							{Timestamp: time.Time{}, Value: 5},
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
			nil, // store batch OK
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
		}

		assert.Equal(t, expected, actual)
	})
}
