package dashboardadmin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/onix/domain"
	mockRepository "github.com/goforbroke1006/onix/mocks/repository"
	"github.com/goforbroke1006/onix/pkg/log"
)

func TestNewServer(t *testing.T) {
	type args struct {
		serviceRepo  domain.ServiceRepository
		releaseRepo  domain.ReleaseRepository
		sourceRepo   domain.SourceRepository
		criteriaRepo domain.CriteriaRepository
		logger       log.Logger
	}
	tests := []struct {
		name string
		args args
		want *server
	}{
		{
			name: "positive - returns instance",
			args: args{
				serviceRepo:  nil,
				releaseRepo:  nil,
				sourceRepo:   nil,
				criteriaRepo: nil,
				logger:       nil,
			},
			want: &server{
				serviceRepo:  nil,
				releaseRepo:  nil,
				sourceRepo:   nil,
				criteriaRepo: nil,
				logger:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.serviceRepo, tt.args.releaseRepo, tt.args.sourceRepo, tt.args.criteriaRepo, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
					repo := mockRepository.NewMockServiceRepository(ctrl)
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
					repo := mockRepository.NewMockServiceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Service{
						{Title: "service 1"},
					}, nil)
					return repo
				},
				releaseRepo: func(ctrl *gomock.Controller) domain.ReleaseRepository {
					repo := mockRepository.NewMockReleaseRepository(ctrl)
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
					repository := mockRepository.NewMockServiceRepository(ctrl)
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
					repo := mockRepository.NewMockServiceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Service{
						{Title: "service 1"},
						{Title: "service 2"},
						{Title: "service 3"},
					}, nil)
					return repo
				},
				releaseRepo: func(ctrl *gomock.Controller) domain.ReleaseRepository {
					repo := mockRepository.NewMockReleaseRepository(ctrl)
					releases := []domain.Release{
						{ID: 1, Service: "service", Name: "1.111.0", StartAt: time.Time{}},
						{ID: 2, Service: "service", Name: "1.112.0", StartAt: time.Time{}},
						{ID: 3, Service: "service", Name: "1.113.0", StartAt: time.Time{}},
					}
					repo.EXPECT().GetNLasts(gomock.Any(), gomock.Any()).Return(releases, nil).Times(3)
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
					repo := mockRepository.NewMockSourceRepository(ctrl)
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
					repo := mockRepository.NewMockSourceRepository(ctrl)
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
			wantStatusCode: http.StatusOK,
		},
		{
			name: "positive 2 - several sources",
			fields: fields{
				sourceRepo: func(ctrl *gomock.Controller) domain.SourceRepository {
					repo := mockRepository.NewMockSourceRepository(ctrl)
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

func Test_server_PostCriteria(t *testing.T) {
	type fields struct {
		criteriaRepo func(ctrl *gomock.Controller) domain.CriteriaRepository
	}
	type args struct {
		postBody string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "negative 1 - DB fail",
			fields: fields{
				criteriaRepo: func(ctrl *gomock.Controller) domain.CriteriaRepository {
					repo := mockRepository.NewMockCriteriaRepository(ctrl)
					repo.EXPECT().
						Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(0), errors.New("tcp fake error"))
					return repo
				},
			},
			args: args{
				postBody: `{
	"service_name": "foo/backend", 
	"title":        "some criteria", 
	"selector":     "some query", 
	"interval":     "5m", 
	"expected_dir": "equal" 
}`,
			},
			wantErr:        true,
			wantStatusCode: 0,
		},
		{
			name:   "negative 1 - broken post data",
			fields: fields{},
			args: args{
				postBody: `{
	"service_name": "foo/backend", 
	"title":        "some criteria", 
	"selector":     "some query", 
	"interval":     "5m", 
	"expected_dir": "eq`,
			},
			wantErr:        true,
			wantStatusCode: 0,
		},
		{
			name: "positive 1 - stored OK",
			fields: fields{
				criteriaRepo: func(ctrl *gomock.Controller) domain.CriteriaRepository {
					repo := mockRepository.NewMockCriteriaRepository(ctrl)
					repo.EXPECT().
						Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(1), nil)
					return repo
				},
			},
			args: args{
				postBody: `{
	"service_name": "foo/backend", 
	"title":        "some criteria", 
	"selector":     "some query", 
	"interval":     "5m", 
	"expected_dir": "equal" 
}`,
			},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var (
				criteriaRepo domain.CriteriaRepository
			)

			if tt.fields.criteriaRepo != nil {
				criteriaRepo = tt.fields.criteriaRepo(mockCtrl)
			}

			s := server{
				criteriaRepo: criteriaRepo,
				logger:       log.NewLogger(),
			}

			req, _ := http.NewRequest(http.MethodGet, "", strings.NewReader(tt.args.postBody))
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
			ctx := echo.New().NewContext(req, rec)

			if err := s.PostCriteria(ctx); (err != nil) != tt.wantErr {
				t.Errorf("PostCriteria() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantStatusCode {
				t.Errorf("GetService() code = %v, want %v", rec.Code, tt.wantStatusCode)
			}
		})
	}
}
