package dashboardmain // nolint:testpackage

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	apiSpec "github.com/goforbroke1006/onix/api/dashboard-main"
	"github.com/goforbroke1006/onix/common"
	"github.com/goforbroke1006/onix/internal/repository"
	"github.com/goforbroke1006/onix/internal/service"
	"github.com/goforbroke1006/onix/tests"
)

func Test_handlers_GetCompare(t *testing.T) { // nolint:funlen,paralleltest
	connString := common.GetTestConnectionStrings()

	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		t.Skip(err)
	}
	defer conn.Close()

	var (
		releaseRepository     = repository.NewReleaseRepository(conn)
		measurementRepository = repository.NewMeasurementRepository(conn)
	)

	handlersInstance := handlers{
		serviceRepo:        nil,
		releaseSvc:         service.NewReleaseService(releaseRepository),
		sourceRepo:         repository.NewSourceRepository(conn),
		criteriaRepo:       repository.NewCriteriaRepository(conn),
		measurementService: service.NewMeasurementService(measurementRepository),
		logger:             nil,
	}

	t.Run("release 1 not found", func(t *testing.T) { // nolint:paralleltest
		const (
			fakeServiceName = "foo/bar/backend"
			sourceID        = 1
			fixture         = "./testdata/handlers_GetCompare_release-1-not-found.fixture.sql"
		)

		if err := tests.LoadFixture(conn, fixture); err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		rr := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, rr)
		err = handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    "2.0.0",
			ReleaseOneStart:    1642877700,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    "2.1.0",
			ReleaseTwoStart:    1643894976,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		actualErr := errors.New("can't get release one by name")
		assert.ErrorAs(t, err, &actualErr)
	})

	t.Run("release 2 not found", func(t *testing.T) { // nolint:paralleltest
		const (
			fakeServiceName = "foo/bar/backend"
			sourceID        = 1
			fixture         = "./testdata/handlers_GetCompare_release-2-not-found.fixture.sql"
		)

		if err := tests.LoadFixture(conn, fixture); err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		rr := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, rr)
		err = handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    "2.0.0",
			ReleaseOneStart:    1642877700,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    "2.1.0",
			ReleaseTwoStart:    1643894976,
			ReleaseTwoSourceId: sourceID,
			Period:             "1h",
		})
		assert.NotNil(t, err)
		actualErr := errors.New("can't get release two by name")
		assert.ErrorAs(t, err, &actualErr)
	})

	t.Run("basic", func(t *testing.T) { // nolint:paralleltest
		const (
			fakeServiceName = "foo/bar/backend"
			sourceID        = 1
			fixture         = "./testdata/handlers_GetCompare_basic.fixture.sql"
		)

		if err := tests.LoadFixture(conn, fixture); err != nil {
			t.Fatal(err)
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "", nil)
		recorder := httptest.NewRecorder()
		echoContext := echo.New().NewContext(req, recorder)
		err = handlersInstance.GetCompare(echoContext, apiSpec.GetCompareParams{
			Service:            fakeServiceName,
			ReleaseOneTitle:    "2.0.0",
			ReleaseOneStart:    1642877700,
			ReleaseOneSourceId: sourceID,
			ReleaseTwoTitle:    "2.1.0",
			ReleaseTwoStart:    1643894976,
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
