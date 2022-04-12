package service // nolint:testpackage

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/internal/repository/mocks"
	"github.com/goforbroke1006/onix/tests"
)

func TestNewReleaseService(t *testing.T) {
	t.Parallel()

	service := NewReleaseService(nil)
	assert.NotNil(t, service)
}

func Test_releaseService_GetReleases(t *testing.T) { // nolint:funlen
	t.Parallel()

	const (
		fakeServiceName = "wildfowl/hello"
	)

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetReleases(gomock.Eq(fakeServiceName), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("db broken"))

		svc := releaseService{repo: releaseRepository}
		_, err := svc.GetReleases(fakeServiceName, time.Time{}, time.Time{})
		assert.NotNil(t, err)
		assert.Equal(t, "can't get releases: db broken", err.Error())
	})

	t.Run("no error but releases list empty", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetReleases(gomock.Eq(fakeServiceName), gomock.Any(), gomock.Any()).Return(nil, nil)

		svc := releaseService{repo: releaseRepository}
		releases, err := svc.GetReleases(fakeServiceName, time.Time{}, time.Time{})
		assert.Nil(t, releases)
		assert.Nil(t, err)
	})

	t.Run("releases list NOT empty", func(t *testing.T) {
		t.Parallel()

		fakeReleases := []domain.Release{
			{ID: 1, Service: fakeServiceName, Name: "v1.0.0", StartAt: time.Time{}},
			{ID: 2, Service: fakeServiceName, Name: "v1.0.1", StartAt: time.Time{}},
			{ID: 3, Service: fakeServiceName, Name: "v1.0.2", StartAt: time.Time{}},
		}

		t.Run("err on getting next", func(t *testing.T) {
			releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
			releaseRepository.EXPECT().GetReleases(gomock.Eq(fakeServiceName), gomock.Any(), gomock.Any()).
				Return(fakeReleases, nil)

			releaseRepository.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Any()).
				Return(nil, errors.New("fake error"))

			svc := releaseService{repo: releaseRepository}
			releases, err := svc.GetReleases(fakeServiceName, time.Time{}, time.Time{})
			assert.Nil(t, releases)
			assert.NotNil(t, err)
			assert.Equal(t, "can't get next release: fake error", err.Error())
		})

		t.Run("no next", func(t *testing.T) {
			releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
			releaseRepository.EXPECT().GetReleases(gomock.Eq(fakeServiceName), gomock.Any(), gomock.Any()).
				Return(fakeReleases, nil)

			releaseRepository.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Any()).Return(nil, nil)

			svc := releaseService{repo: releaseRepository}
			releases, err := svc.GetReleases(fakeServiceName, time.Time{}, time.Time{})
			assert.NotNil(t, releases)
			assert.Equal(t, 3, len(releases))
			assert.Nil(t, err)
		})

		t.Run("has next", func(t *testing.T) {
			releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
			releaseRepository.EXPECT().GetReleases(gomock.Eq(fakeServiceName), gomock.Any(), gomock.Any()).
				Return(fakeReleases, nil)

			releaseRepository.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Any()).
				Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)

			svc := releaseService{repo: releaseRepository}
			releases, err := svc.GetReleases(fakeServiceName, time.Time{}, time.Time{})
			assert.NotNil(t, releases)
			assert.Equal(t, len(fakeReleases), len(releases))
			assert.Nil(t, err)
		})
	})
}

func Test_releaseService_GetAll(t *testing.T) { // nolint:funlen
	t.Parallel()

	type fields struct {
		repo domain.ReleaseRepository
	}

	type args struct {
		serviceName string
	}

	testsCases := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.ReleaseTimeRange
		wantErr bool
	}{
		{
			name:   "basic",
			fields: fields{repo: stubReleaseRepository{}},
			args:   args{serviceName: "foo/bar/backend"},
			want: []domain.ReleaseTimeRange{
				{
					ID:      1,
					Service: "foo/bar/backend",
					Name:    "v1.10.0",
					StartAt: time.Date(2018, time.February, 1, 12, 0, 0, 0, time.UTC),
					StopAt:  time.Date(2018, time.August, 2, 11, 59, 59, 0, time.UTC),
				},
				{
					ID:      2,
					Service: "foo/bar/backend",
					Name:    "v1.11.0",
					StartAt: time.Date(2018, time.August, 2, 12, 0, 0, 0, time.UTC),
					StopAt:  time.Now().Truncate(time.Second).UTC(),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testsCases {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			svc := releaseService{
				repo: testCase.fields.repo,
			}
			got, err := svc.GetAll(testCase.args.serviceName)
			if (err != nil) != testCase.wantErr {
				t.Fatalf("GetAll() error = %v, wantErr %v", err, testCase.wantErr)
			}

			if len(got) != len(testCase.want) {
				t.Errorf("GetAll() got = %v, want %v", len(got), len(testCase.want))
			}

			for releaseIndex := range got {
				if got[releaseIndex].ID != testCase.want[releaseIndex].ID {
					t.Errorf("GetAll() got = %v, want %v", got[releaseIndex].ID, testCase.want[releaseIndex].ID)
				}
				if got[releaseIndex].Service != testCase.want[releaseIndex].Service {
					t.Errorf("GetAll() got = %v, want %v", got[releaseIndex].Service, testCase.want[releaseIndex].Service)
				}
				if got[releaseIndex].Name != testCase.want[releaseIndex].Name {
					t.Errorf("GetAll() got = %v, want %v", got[releaseIndex].Name, testCase.want[releaseIndex].Name)
				}
				if got[releaseIndex].StartAt != testCase.want[releaseIndex].StartAt {
					t.Errorf("GetAll() got = %v, want %v", got[releaseIndex].StartAt, testCase.want[releaseIndex].StartAt)
				}
				if got[releaseIndex].StopAt.Sub(testCase.want[releaseIndex].StopAt) > 5*time.Second {
					t.Errorf("GetAll() got = %v, want %v", got[releaseIndex].StopAt, testCase.want[releaseIndex].StopAt)
				}
			}
		})
	}

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetAll(gomock.All()).Return(nil, errors.New("fake error"))
		svc := releaseService{
			repo: releaseRepository,
		}

		releases, err := svc.GetAll("")
		assert.NotNil(t, err)
		assert.Equal(t, "can't get releases: fake error", err.Error())
		assert.Nil(t, releases)
	})

	t.Run("repo ok but releases list is empty", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetAll(gomock.All()).Return(nil, nil)
		svc := releaseService{
			repo: releaseRepository,
		}

		releases, err := svc.GetAll("")
		assert.Nil(t, err)
		assert.Nil(t, releases)
	})
}

func Test_releaseService_GetByName(t *testing.T) { // nolint:funlen
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("repo error", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(nil, errors.New("fake error"))
		svc := releaseService{
			repo: releaseRepository,
		}

		release, err := svc.GetByName("", "")
		assert.NotNil(t, err)
		assert.Equal(t, "can't get release: fake error", err.Error())
		assert.Nil(t, release)
	})

	t.Run("repo ok but current is NIL", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetByName(gomock.Any(), gomock.Any()).Return(nil, nil)
		svc := releaseService{
			repo: releaseRepository,
		}

		release, err := svc.GetByName("", "")
		assert.NotNil(t, err)
		assert.Equal(t, domain.ErrNotFound, err)
		assert.Nil(t, release)
	})

	t.Run("repo ok but current found but next with err", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetByName(gomock.Any(), gomock.Any()).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepository.EXPECT().GetNextAfter(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("fake error"))
		svc := releaseService{
			repo: releaseRepository,
		}

		release, err := svc.GetByName("", "")
		assert.NotNil(t, err)
		assert.Equal(t, "can't get next release: fake error", err.Error())
		assert.Nil(t, release)
	})

	t.Run("repo ok but current found but next NIL", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetByName(gomock.Any(), gomock.Any()).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepository.EXPECT().GetNextAfter(gomock.Any(), gomock.Any()).Return(nil, nil)
		svc := releaseService{
			repo: releaseRepository,
		}

		release, err := svc.GetByName("", "")
		assert.Nil(t, err)
		assert.NotNil(t, release)
	})

	t.Run("repo ok but current found and next found", func(t *testing.T) {
		t.Parallel()

		releaseRepository := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepository.EXPECT().GetByName(gomock.Any(), gomock.Any()).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepository.EXPECT().GetNextAfter(gomock.Any(), gomock.Any()).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		svc := releaseService{
			repo: releaseRepository,
		}

		release, err := svc.GetByName("", "")
		assert.Nil(t, err)
		assert.NotNil(t, release)
	})
}

var _ domain.ReleaseRepository = &stubReleaseRepository{}

type stubReleaseRepository struct{}

func (repo stubReleaseRepository) Store(serviceName string, releaseName string, startAt time.Time) error {
	panic("implement me")
}

func (repo stubReleaseRepository) GetReleases(serviceName string, from, till time.Time) ([]domain.Release, error) {
	panic("implement me")
}

func (repo stubReleaseRepository) GetByName(serviceName, releaseName string) (*domain.Release, error) {
	panic("implement me")
}

func (repo stubReleaseRepository) GetNextAfter(serviceName, releaseName string) (*domain.Release, error) {
	panic("implement me")
}

func (repo stubReleaseRepository) GetLast(serviceName string) (*domain.Release, error) {
	panic("implement me")
}

func (repo stubReleaseRepository) GetNLasts(serviceName string, count uint) ([]domain.Release, error) {
	panic("implement me")
}

func (repo stubReleaseRepository) GetAll(serviceName string) ([]domain.Release, error) {
	var releases []domain.Release

	if serviceName == "foo/bar/backend" {
		releases = []domain.Release{
			{
				ID:      1,
				Service: "foo/bar/backend",
				Name:    "v1.10.0",
				StartAt: time.Date(2018, time.February, 1, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:      2,
				Service: "foo/bar/backend",
				Name:    "v1.11.0",
				StartAt: time.Date(2018, time.August, 2, 12, 0, 0, 0, time.UTC),
			},
		}
	}

	return releases, nil
}

func TestGetReleases(t *testing.T) { // nolint:paralleltest
	connString := common.GetTestConnectionStrings()

	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		t.Skip(err)
	}
	defer conn.Close()

	releaseRepository := repository.NewReleaseRepository(conn)
	releaseService := NewReleaseService(releaseRepository)

	t.Run("inside range", func(t *testing.T) { // nolint:paralleltest
		if err := tests.LoadFixture(conn, "./testdata/release_test.fixture.sql"); err != nil {
			t.Fatal(err)
		}

		from, _ := time.Parse("2006-01-02 15:04:05", "2020-10-25 00:00:00")
		till, _ := time.Parse("2006-01-02 15:04:05", "2020-11-06 00:00:00")
		ranges, err := releaseService.GetReleases("foo/bar/backend", from, till)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(ranges))
		assert.Equal(t, time.Date(2020, time.October, 25, 0, 0, 0, 0, time.UTC), ranges[0].StartAt)
		assert.Equal(t, time.Date(2020, time.October, 25, 23, 59, 59, 0, time.UTC), ranges[0].StopAt)

		assert.Equal(t, time.Date(2020, time.November, 6, 0, 0, 0, 0, time.UTC), ranges[2].StartAt)
		assert.Equal(t, time.Date(2020, time.November, 13, 23, 59, 59, 0, time.UTC), ranges[2].StopAt)
	})

	t.Run("in the end of range", func(t *testing.T) { // nolint:paralleltest
		if err := tests.LoadFixture(conn, "./testdata/release_test.fixture.sql"); err != nil {
			t.Fatal(err)
		}

		from, _ := time.Parse("2006-01-02 15:04:05", "2020-11-28 00:00:00")
		till, _ := time.Parse("2006-01-02 15:04:05", "2020-12-26 00:00:00")
		ranges, err := releaseService.GetReleases("foo/bar/backend", from, till)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(ranges))

		assert.Equal(t, time.Date(2020, time.December, 26, 0, 0, 0, 0, time.UTC), ranges[2].StartAt)
		assert.Less(t, time.Now().UTC().Sub(ranges[2].StopAt), 2*time.Second)
	})
}
