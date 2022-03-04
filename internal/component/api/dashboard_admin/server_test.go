package dashboard_admin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/repository/mocks"
)

func Test_server_GetService(t *testing.T) {
	type fields struct {
		serviceRepo func(ctrl *gomock.Controller) domain.ServiceRepository
		releaseRepo func(ctrl *gomock.Controller) domain.ReleaseRepository
	}
	type args struct{}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "negative 1 - DB fail on load services",
			fields: fields{
				serviceRepo: func(ctrl *gomock.Controller) domain.ServiceRepository {
					repo := mocks.NewMockServiceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Service{}, errors.New("tcp lost connection"))
					return repo
				},
			},
			args:           args{},
			wantErr:        true,
			wantStatusCode: 0,
		},
		{
			name: "negative 2 - DB fail on load releases",
			fields: fields{
				serviceRepo: func(ctrl *gomock.Controller) domain.ServiceRepository {
					repo := mocks.NewMockServiceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Service{
						{Title: "service 1"},
					}, nil)
					return repo
				},
				releaseRepo: func(ctrl *gomock.Controller) domain.ReleaseRepository {
					repo := mocks.NewMockReleaseRepository(ctrl)
					repo.EXPECT().GetNLasts(gomock.Any(), gomock.Any()).Return([]domain.Release{}, errors.New("tcp lost connection"))
					return repo
				},
			},
			args:           args{},
			wantErr:        true,
			wantStatusCode: 0,
		},
		{
			name: "positive 1 - no services in DB",
			fields: fields{
				serviceRepo: func(ctrl *gomock.Controller) domain.ServiceRepository {
					repository := mocks.NewMockServiceRepository(ctrl)
					repository.EXPECT().GetAll().Return([]domain.Service{}, nil)
					return repository
				},
			},
			args:           args{},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "positive 1 - has services in DB",
			fields: fields{
				serviceRepo: func(ctrl *gomock.Controller) domain.ServiceRepository {
					repo := mocks.NewMockServiceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Service{
						{Title: "service 1"},
						{Title: "service 2"},
						{Title: "service 3"},
					}, nil)
					return repo
				},
				releaseRepo: func(ctrl *gomock.Controller) domain.ReleaseRepository {
					repo := mocks.NewMockReleaseRepository(ctrl)
					releases := []domain.Release{
						{ID: 1, Service: "service", Name: "1.111.0", StartAt: time.Time{}},
						{ID: 2, Service: "service", Name: "1.112.0", StartAt: time.Time{}},
						{ID: 3, Service: "service", Name: "1.113.0", StartAt: time.Time{}},
					}
					repo.EXPECT().GetNLasts(gomock.Any(), gomock.Any()).Return(releases, nil)
					repo.EXPECT().GetNLasts(gomock.Any(), gomock.Any()).Return(releases, nil)
					repo.EXPECT().GetNLasts(gomock.Any(), gomock.Any()).Return(releases, nil)
					return repo
				},
			},
			args:           args{},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var (
				serviceRepo domain.ServiceRepository
				releaseRepo domain.ReleaseRepository
			)

			if tt.fields.serviceRepo != nil {
				serviceRepo = tt.fields.serviceRepo(mockCtrl)
			}
			if tt.fields.releaseRepo != nil {
				releaseRepo = tt.fields.releaseRepo(mockCtrl)
			}

			s := server{
				serviceRepo: serviceRepo,
				releaseRepo: releaseRepo,
			}

			req, _ := http.NewRequest(http.MethodGet, "", nil)
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
			ctx := echo.New().NewContext(req, rec)

			if err := s.GetService(ctx); (err != nil) != tt.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantStatusCode {
				t.Errorf("GetService() code = %v, want %v", rec.Code, tt.wantStatusCode)
			}
		})
	}
}

func Test_server_GetSource(t *testing.T) {
	type fields struct {
		sourceRepo func(ctrl *gomock.Controller) domain.SourceRepository
	}
	type args struct{}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "negative 1 - DB failed",
			fields: fields{
				sourceRepo: func(ctrl *gomock.Controller) domain.SourceRepository {
					repo := mocks.NewMockSourceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Source{}, errors.New("tcp lost connection"))
					return repo
				},
			},
			args:           args{},
			wantErr:        true,
			wantStatusCode: 0,
		},
		{
			name: "positive 1 - no source in DB",
			fields: fields{
				sourceRepo: func(ctrl *gomock.Controller) domain.SourceRepository {
					repo := mocks.NewMockSourceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Source{
						{ID: 1, Title: "source 1", Kind: "", Address: ""},
						{ID: 2, Title: "source 2", Kind: "", Address: ""},
						{ID: 3, Title: "source 3", Kind: "", Address: ""},
					}, nil)
					return repo
				},
			},
			args:           args{},
			wantErr:        false,
			wantStatusCode: 200,
		},
		{
			name: "positive 2 - several sources",
			fields: fields{
				sourceRepo: func(ctrl *gomock.Controller) domain.SourceRepository {
					repo := mocks.NewMockSourceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Source{
						{ID: 1, Title: "source 1", Kind: "", Address: ""},
						{ID: 2, Title: "source 2", Kind: "", Address: ""},
						{ID: 3, Title: "source 3", Kind: "", Address: ""},
					}, nil)
					return repo
				},
			},
			args:           args{},
			wantErr:        false,
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var (
				sourceRepo domain.SourceRepository
			)

			if tt.fields.sourceRepo != nil {
				sourceRepo = tt.fields.sourceRepo(mockCtrl)
			}

			s := server{
				sourceRepo: sourceRepo,
			}

			req, _ := http.NewRequest(http.MethodGet, "", nil)
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
			ctx := echo.New().NewContext(req, rec)

			if err := s.GetSource(ctx); (err != nil) != tt.wantErr {
				t.Errorf("GetSource() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantStatusCode {
				t.Errorf("GetService() code = %v, want %v", rec.Code, tt.wantStatusCode)
			}
		})
	}
}
