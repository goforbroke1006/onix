package service

import (
	"context"
	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/tests"
	"github.com/jackc/pgx/v4/pgxpool"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/domain"
)

func TestNewReleaseService(t *testing.T) {
	service := NewReleaseService(nil)
	assert.NotNil(t, service)
}

func Test_releaseService_GetAll(t *testing.T) {
	type fields struct {
		repo domain.ReleaseRepository
	}
	type args struct {
		serviceName string
	}
	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := releaseService{
				repo: tt.fields.repo,
			}
			got, err := svc.GetAll(tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("GetAll() got = %v, want %v", len(got), len(tt.want))
			}

			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("GetAll() got = %v, want %v", got[i].ID, tt.want[i].ID)
				}
				if got[i].Service != tt.want[i].Service {
					t.Errorf("GetAll() got = %v, want %v", got[i].Service, tt.want[i].Service)
				}
				if got[i].Name != tt.want[i].Name {
					t.Errorf("GetAll() got = %v, want %v", got[i].Name, tt.want[i].Name)
				}
				if got[i].StartAt != tt.want[i].StartAt {
					t.Errorf("GetAll() got = %v, want %v", got[i].StartAt, tt.want[i].StartAt)
				}
				if got[i].StopAt.Sub(tt.want[i].StopAt) > 5*time.Second {
					t.Errorf("GetAll() got = %v, want %v", got[i].StopAt, tt.want[i].StopAt)
				}
			}
		})
	}
}

var _ domain.ReleaseRepository = &stubReleaseRepository{}

type stubReleaseRepository struct{}

func (repo stubReleaseRepository) Store(serviceName string, releaseName string, startAt time.Time) error {
	// TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetReleases(serviceName string, from, till time.Time) ([]domain.Release, error) {
	// TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetByName(serviceName, releaseName string) (*domain.Release, error) {
	// TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetNextAfter(serviceName, releaseName string) (*domain.Release, error) {
	// TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetLast(serviceName string) (*domain.Release, error) {
	// TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetNLasts(serviceName string, count uint) ([]domain.Release, error) {
	// TODO implement me
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

func TestGetReleases(t *testing.T) {
	connString := common.GetTestConnectionStrings()
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		t.Skip(err)
	}
	defer conn.Close()

	releaseRepository := repository.NewReleaseRepository(conn)
	releaseService := NewReleaseService(releaseRepository)

	t.Run("inside range", func(t *testing.T) {
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

	t.Run("in the end of range", func(t *testing.T) {
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
