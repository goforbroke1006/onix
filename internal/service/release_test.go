//go:build unit
// +build unit

package service

import (
	"testing"
	"time"

	"github.com/goforbroke1006/onix/domain"
)

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

var (
	_ domain.ReleaseRepository = &stubReleaseRepository{}
)

type stubReleaseRepository struct {
}

func (repo stubReleaseRepository) Store(serviceName string, releaseName string, startAt time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetReleases(serviceName string, from, till time.Time) ([]domain.Release, error) {
	//TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetByName(serviceName, releaseName string) (*domain.Release, error) {
	//TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetNextAfter(serviceName, releaseName string) (*domain.Release, error) {
	//TODO implement me
	panic("implement me")
}

func (repo stubReleaseRepository) GetLast(serviceName string) (*domain.Release, error) {
	//TODO implement me
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
