package service

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
	actual := NewMeasurementService(nil)
	assert.NotNil(t, actual)
}

func Test_measurementService_getTimePoints(t *testing.T) {
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
				time.Date(2022, time.March, 8, 14, 00, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := measurementService{}
			if got := m.getTimePoints(tt.args.from, tt.args.till, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTimePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeasurementService_GetOrPull(t *testing.T) {
	t.Parallel()

	t.Run("on DB failure", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		measurementRepository := mockRepo.NewMockMeasurementRepository(mockCtrl)
		svc := measurementService{
			measurementRepo: measurementRepository,
		}

		measurementRepository.EXPECT().GetForPoints(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nil,
			errors.New("some fake error"),
		)

		source := domain.Source{
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{}

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
			measurementRepo: measurementRepository,
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

		source := domain.Source{
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{}

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

		source := domain.Source{
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{}

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

		source := domain.Source{
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{}

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

		source := domain.Source{
			Type: domain.SourceTypePrometheus,
		}
		criteria := domain.Criteria{}

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
