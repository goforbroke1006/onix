package dashboardadmin // nolint:testpackage

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/goforbroke1006/onix/domain"
	mockRepository "github.com/goforbroke1006/onix/internal/repository/mocks"
	"github.com/goforbroke1006/onix/pkg/log"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

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
		want *handlers
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
			want: &handlers{
				serviceRepo:  nil,
				releaseRepo:  nil,
				sourceRepo:   nil,
				criteriaRepo: nil,
				logger:       nil,
			},
		},
	}
	for _, tt := range tests {
		ttCase := tt
		t.Run(ttCase.name, func(t *testing.T) {
			t.Parallel()

			if got := NewHandlers(
				ttCase.args.serviceRepo,
				ttCase.args.releaseRepo,
				ttCase.args.sourceRepo,
				ttCase.args.criteriaRepo,
				ttCase.args.logger,
			); !reflect.DeepEqual(got, ttCase.want) {
				t.Errorf("NewHandlers() = %v, want %v", got, ttCase.want)
			}
		})
	}
}

func Test_handlers_GetHealthz(t *testing.T) {
	t.Parallel()

	var target handlers

	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
	recorder := httptest.NewRecorder()
	echoContext := echo.New().NewContext(req, recorder)

	err := target.GetHealthz(echoContext)
	assert.Nil(t, err)
}

func Test_handlers_GetService(t *testing.T) { // nolint:funlen
	t.Parallel()

	type fields struct {
		serviceRepo func(ctrl *gomock.Controller) domain.ServiceRepository
		releaseRepo func(ctrl *gomock.Controller) domain.ReleaseRepository
	}

	tests := []struct {
		name           string
		fields         fields
		wantErr        bool
		wantStatusCode int
	}{
		{
			name: "negative 1 - DB fail on load services",
			fields: fields{
				serviceRepo: func(ctrl *gomock.Controller) domain.ServiceRepository {
					repo := mockRepository.NewMockServiceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Service{}, errors.New("tcp lost connection")) // nolint:goerr113

					return repo
				},
				releaseRepo: nil,
			},
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
					repo.EXPECT().GetNLasts(gomock.Any(), gomock.Any()).Return(
						[]domain.Release{},
						errors.New("tcp lost connection"), // nolint:goerr113
					)

					return repo
				},
			},
			wantErr:        true,
			wantStatusCode: 0,
		},
		{
			name: "positive 1 - no services in DB",
			fields: fields{
				serviceRepo: func(ctrl *gomock.Controller) domain.ServiceRepository {
					repo := mockRepository.NewMockServiceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Service{}, nil)

					return repo
				},
				releaseRepo: nil,
			},
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
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		ttCase := tt
		t.Run(ttCase.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var (
				serviceRepo domain.ServiceRepository
				releaseRepo domain.ReleaseRepository
			)

			if ttCase.fields.serviceRepo != nil {
				serviceRepo = ttCase.fields.serviceRepo(mockCtrl)
			}
			if ttCase.fields.releaseRepo != nil {
				releaseRepo = ttCase.fields.releaseRepo(mockCtrl)
			}

			instance := handlers{
				serviceRepo:  serviceRepo,
				releaseRepo:  releaseRepo,
				sourceRepo:   nil,
				criteriaRepo: nil,
				logger:       nil,
			}

			req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})} // nolint:exhaustivestruct
			ctx := echo.New().NewContext(req, rec)

			if err := instance.GetService(ctx); (err != nil) != ttCase.wantErr {
				t.Errorf("GetService() error = %v, wantErr %v", err, ttCase.wantErr)
			}

			if rec.Code != ttCase.wantStatusCode {
				t.Errorf("GetService() code = %v, want %v", rec.Code, ttCase.wantStatusCode)
			}
		})
	}
}

func Test_handlers_GetSource(t *testing.T) { // nolint:funlen
	t.Parallel()

	type fields struct {
		sourceRepo func(ctrl *gomock.Controller) domain.SourceRepository
	}

	type testCase struct {
		name           string
		fields         fields
		wantErr        bool
		wantStatusCode int
	}

	tests := []testCase{
		{
			name: "negative 1 - DB failed",
			fields: fields{
				sourceRepo: func(ctrl *gomock.Controller) domain.SourceRepository {
					repo := mockRepository.NewMockSourceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Source{}, errors.New("tcp lost connection")) // nolint:goerr113

					return repo
				},
			},
			wantErr:        true,
			wantStatusCode: 0,
		},
		{
			name: "positive 1 - no source in DB",
			fields: fields{
				sourceRepo: func(ctrl *gomock.Controller) domain.SourceRepository {
					repo := mockRepository.NewMockSourceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Source{
						{ID: 1, Title: "source 1", Type: "", Address: ""},
						{ID: 2, Title: "source 2", Type: "", Address: ""},
						{ID: 3, Title: "source 3", Type: "", Address: ""},
					}, nil)

					return repo
				},
			},
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "positive 2 - several sources",
			fields: fields{
				sourceRepo: func(ctrl *gomock.Controller) domain.SourceRepository {
					repo := mockRepository.NewMockSourceRepository(ctrl)
					repo.EXPECT().GetAll().Return([]domain.Source{
						{ID: 1, Title: "source 1", Type: "", Address: ""},
						{ID: 2, Title: "source 2", Type: "", Address: ""},
						{ID: 3, Title: "source 3", Type: "", Address: ""},
					}, nil)

					return repo
				},
			},
			wantErr:        false,
			wantStatusCode: 200,
		},
	}

	for _, tt := range tests {
		ttCase := tt
		func(ttCase testCase) {
			t.Run(ttCase.name, func(t *testing.T) {
				t.Parallel()

				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				var sourceRepo domain.SourceRepository

				if ttCase.fields.sourceRepo != nil {
					sourceRepo = ttCase.fields.sourceRepo(mockCtrl)
				}

				instance := handlers{
					serviceRepo:  nil,
					releaseRepo:  nil,
					sourceRepo:   sourceRepo,
					criteriaRepo: nil,
					logger:       nil,
				}

				req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
				rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})} // nolint:exhaustivestruct
				ctx := echo.New().NewContext(req, rec)

				if err := instance.GetSource(ctx); (err != nil) != ttCase.wantErr {
					t.Errorf("GetSource() error = %v, wantErr %v", err, ttCase.wantErr)
				}

				if rec.Code != ttCase.wantStatusCode {
					t.Errorf("GetService() code = %v, want %v", rec.Code, ttCase.wantStatusCode)
				}
			})
		}(ttCase)
	}
}

func Test_handlers_PostCriteria(t *testing.T) { // nolint:funlen
	t.Parallel()

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
						Return(int64(0), errors.New("tcp fake error")) // nolint:goerr113

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
			name: "negative 1 - broken post data",
			fields: fields{
				criteriaRepo: nil,
			},
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
		ttCase := tt
		t.Run(ttCase.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var criteriaRepo domain.CriteriaRepository

			if ttCase.fields.criteriaRepo != nil {
				criteriaRepo = ttCase.fields.criteriaRepo(mockCtrl)
			}

			instance := handlers{
				serviceRepo:  nil,
				releaseRepo:  nil,
				sourceRepo:   nil,
				criteriaRepo: criteriaRepo,
				logger:       log.NewNullLogger(),
			}

			req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", strings.NewReader(ttCase.args.postBody))
			rec := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})} // nolint:exhaustivestruct
			ctx := echo.New().NewContext(req, rec)

			if err := instance.PostCriteria(ctx); (err != nil) != ttCase.wantErr {
				t.Errorf("PostCriteria() error = %v, wantErr %v", err, ttCase.wantErr)
			}

			if rec.Code != ttCase.wantStatusCode {
				t.Errorf("GetService() code = %v, want %v", rec.Code, ttCase.wantStatusCode)
			}
		})
	}
}
