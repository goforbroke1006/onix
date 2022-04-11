package dashboardmain // nolint:testpackage

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-main"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/internal/repository/mocks"
	"github.com/goforbroke1006/onix/internal/service"
	"github.com/goforbroke1006/onix/pkg/log"
)

func TestNewHandlers(t *testing.T) {
	t.Parallel()

	h := NewHandlers(nil, nil, nil, nil, nil, nil)
	assert.NotNil(t, h)
}

func TestHandlers_GetHealthz(t *testing.T) {
	t.Parallel()

	var handlersInstance handlers

	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
	recorder := httptest.NewRecorder()
	echoContext := echo.New().NewContext(req, recorder)

	err := handlersInstance.GetHealthz(echoContext)
	assert.Nil(t, err)
}

func TestHandlers_GetService(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("on db error", func(t *testing.T) {
		t.Parallel()

		serviceRepo := mocks.NewMockServiceRepository(mockCtrl)
		serviceRepo.EXPECT().GetAll().Return(nil, errors.New("fake db error"))

		handlersInstance := handlers{
			serviceRepo:        serviceRepo,
			releaseSvc:         nil,
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetService(echoContext)
		assert.NotNil(t, err)
	})

	t.Run("basic usage", func(t *testing.T) {
		t.Parallel()

		serviceRepo := mocks.NewMockServiceRepository(mockCtrl)
		serviceRepo.EXPECT().GetAll().Return([]domain.Service{{}, {}, {}}, nil)

		handlersInstance := handlers{
			serviceRepo:        serviceRepo,
			releaseSvc:         nil,
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetService(echoContext)
		assert.Nil(t, err)

		var actual []apiSpec.Service
		respBytes, _ := ioutil.ReadAll(recorder.Body)
		if err := json.Unmarshal(respBytes, &actual); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 3, len(actual))
	})
}

func TestHandlers_GetSource(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("on db error", func(t *testing.T) {
		t.Parallel()

		sourceRepo := mocks.NewMockSourceRepository(mockCtrl)
		sourceRepo.EXPECT().GetAll().Return(nil, errors.New("fake db error"))

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         nil,
			sourceRepo:         sourceRepo,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetSource(echoContext)
		assert.NotNil(t, err)
	})

	t.Run("basic usage", func(t *testing.T) {
		t.Parallel()

		sourceRepo := mocks.NewMockSourceRepository(mockCtrl)
		sourceRepo.EXPECT().GetAll().Return([]domain.Source{{}, {}, {}, {}}, nil)

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         nil,
			sourceRepo:         sourceRepo,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetSource(echoContext)
		assert.Nil(t, err)

		var actual []apiSpec.Source
		respBytes, _ := ioutil.ReadAll(recorder.Body)
		if err := json.Unmarshal(respBytes, &actual); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 4, len(actual))
	})
}

func TestHandlers_GetRelease(t *testing.T) { // nolint:funlen
	t.Parallel()

	const fakeServiceName = "foo/bar/backend"

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("on repo error", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetAll(gomock.Eq(fakeServiceName)).
			Return(nil, errors.New("fake db error"))

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetRelease(echoContext, apiSpec.GetReleaseParams{
			Service: fakeServiceName,
		})
		assert.NotNil(t, err)
	})

	t.Run("basic usage", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetAll(gomock.Eq(fakeServiceName)).
			Return([]domain.Release{
				{ID: 2, Service: fakeServiceName, Name: "v3.4.5", StartAt: time.Unix(1000, 0)},
				{ID: 3, Service: fakeServiceName, Name: "v3.4.7", StartAt: time.Unix(2000, 0)},
			}, nil)

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetRelease(echoContext, apiSpec.GetReleaseParams{
			Service: fakeServiceName,
		})
		assert.Nil(t, err)

		var actual []apiSpec.Release
		respBytes, _ := ioutil.ReadAll(recorder.Body)
		if err := json.Unmarshal(respBytes, &actual); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(2), actual[0].Id)
		assert.Equal(t, "v3.4.5", actual[0].Title)
		assert.Equal(t, int64(1000), actual[0].From)
		assert.Equal(t, int64(1999), actual[0].Till)

		assert.Equal(t, int64(3), actual[1].Id)
		assert.Equal(t, "v3.4.7", actual[1].Title)
		assert.Equal(t, int64(2000), actual[1].From)
		assert.True(t, actual[0].Till-time.Now().Unix() < 5)
	})
}

func Test_handlers_GetCompare(t *testing.T) { // nolint:funlen,maintidx
	t.Parallel()

	const (
		fakeServiceName = "foo/bar/backend"
		sourceID        = int64(1)
		releaseOne      = "2.0.0"
		releaseOneStart = int64(1642877700) // Sat Jan 22 2022 18:55:00 GMT+0000
		releaseTwo      = "2.1.0"
		releaseTwoStart = int64(1643894976) // Thu Feb 03 2022 13:29:36 GMT+0000
	)

	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})

	t.Run("release 1 not found", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(nil, errors.New("fake not found")) // simulation here

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		rr := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, rr)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    1642877700,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    1643894976,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "can't get release one by name: can't get release: fake not found", err.Error())
	})

	t.Run("ask time before release 1 starts", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Unix(1000, 0)}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Unix(2000, 0)}, nil)

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             log.NewNullLogger(),
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		rr := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, rr)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    999,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    1643894976,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "release 1 wrong time: wrong time range", err.Error())
	})

	t.Run("release 2 not found", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(nil, errors.New("fake not found")) // simulation here

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		rr := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, rr)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    1642877700,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    1643894976,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "can't get release two by name: can't get release: fake not found", err.Error())
	})

	t.Run("ask time before release 2 starts", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Unix(1000, 0)}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Unix(2000, 0)}, nil)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Unix(3000, 0)}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Unix(4000, 0)}, nil)

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         nil,
			criteriaRepo:       nil,
			measurementService: nil,
			logger:             log.NewNullLogger(),
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		rr := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, rr)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    1500,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    2999,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "release 2 wrong time: wrong time range", err.Error())
	})

	t.Run("criteria repo error", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)

		criteriaRepo := mocks.NewMockCriteriaRepository(mockCtrl)
		criteriaRepo.EXPECT().GetAll(gomock.Eq(fakeServiceName)).
			Return(nil, errors.New("fake error"))

		handlersInstance := handlers{ // nolint:exhaustivestruct
			releaseSvc:   service.NewReleaseService(releaseRepo),
			criteriaRepo: criteriaRepo,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    releaseOneStart,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    releaseTwoStart,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "can't get criteria list: fake error", err.Error())
	})

	t.Run("source 1 not found", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)

		criteriaRepo := mocks.NewMockCriteriaRepository(mockCtrl)
		criteriaRepo.EXPECT().GetAll(gomock.Eq(fakeServiceName)).
			Return([]domain.Criteria{{}, {}}, nil)

		sourceRepo := mocks.NewMockSourceRepository(mockCtrl)
		sourceRepo.EXPECT().Get(gomock.Eq(sourceID)).
			Return(nil, errors.New("fake error")) // simulation here

		handlersInstance := handlers{ // nolint:exhaustivestruct
			releaseSvc:   service.NewReleaseService(releaseRepo),
			sourceRepo:   sourceRepo,
			criteriaRepo: criteriaRepo,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    releaseOneStart,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    releaseTwoStart,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "can't get source #1: fake error", err.Error())
	})

	t.Run("source 2 not found", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)

		criteriaRepo := mocks.NewMockCriteriaRepository(mockCtrl)
		criteriaRepo.EXPECT().GetAll(gomock.Eq(fakeServiceName)).
			Return([]domain.Criteria{{}, {}}, nil)

		sourceRepo := mocks.NewMockSourceRepository(mockCtrl)
		sourceRepo.EXPECT().Get(gomock.Eq(sourceID)).
			Return(
				&domain.Source{}, // nolint:exhaustivestruct
				nil,
			)
		sourceRepo.EXPECT().Get(gomock.Eq(sourceID)).
			Return(nil, errors.New("fake error")) // simulation here

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         sourceRepo,
			criteriaRepo:       criteriaRepo,
			measurementService: nil,
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    releaseOneStart,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    releaseTwoStart,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "can't get source #2: fake error", err.Error())
	})

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		releaseRepo := mocks.NewMockReleaseRepository(mockCtrl)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseOne)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetByName(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)
		releaseRepo.EXPECT().GetNextAfter(gomock.Eq(fakeServiceName), gomock.Eq(releaseTwo)).
			Return(&domain.Release{ID: 0, Service: "", Name: "", StartAt: time.Time{}}, nil)

		criteriaRepo := mocks.NewMockCriteriaRepository(mockCtrl)
		criteriaRepo.EXPECT().GetAll(gomock.Eq(fakeServiceName)).
			Return([]domain.Criteria{
				{},
				{},
			}, nil)

		sourceRepo := mocks.NewMockSourceRepository(mockCtrl)
		sourceRepo.EXPECT().Get(gomock.Eq(sourceID)).
			Return(&domain.Source{ID: sourceID, Title: "", Type: domain.SourceTypePrometheus, Address: ""}, nil).Times(2)

		measurementRepo := mocks.NewMockMeasurementRepository(mockCtrl)
		measurementRepo.EXPECT().GetForPoints(gomock.Eq(sourceID), gomock.Any(), gomock.Any()).
			Return(make([]domain.MeasurementRow, 12+1), nil).
			Times(4)

		handlersInstance := handlers{
			serviceRepo:        nil,
			releaseSvc:         service.NewReleaseService(releaseRepo),
			sourceRepo:         sourceRepo,
			criteriaRepo:       criteriaRepo,
			measurementService: service.NewMeasurementService(measurementRepo),
			logger:             nil,
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err := handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    releaseOne,
			ReleaseOneStart:    releaseOneStart,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    releaseTwo,
			ReleaseTwoStart:    releaseTwoStart,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.Nil(t, err)

		respBytes, _ := ioutil.ReadAll(recorder.Body)
		var responseObj apiSpec.CompareResponse
		if err := json.Unmarshal(respBytes, &responseObj); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, fakeServiceName, responseObj.Service)
		assert.Equal(t, 2, len(responseObj.Reports))
	})
}
