package system

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	apiSpec "github.com/goforbroke1006/onix/api/system"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewServer creates new server's handlers implementations instance
func NewServer(
	serviceRepo domain.ServiceRepository,
	releaseRepo domain.ReleaseRepository,
	logger log.Logger,
) *server { // nolint:golint
	return &server{
		serviceRepo: serviceRepo,
		releaseRepo: releaseRepo,
		logger:      logger,
	}
}

var (
	_ apiSpec.ServerInterface = &server{}
)

type server struct {
	serviceRepo domain.ServiceRepository
	releaseRepo domain.ReleaseRepository
	logger      log.Logger
}

func (s server) GetHealthz(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (s server) GetRegister(ctx echo.Context, params apiSpec.GetRegisterParams) error {
	startAt := time.Now().UTC()
	if params.StartAt != nil {
		startAt = time.Unix(*params.StartAt, 0).UTC()
	}

	if err := s.serviceRepo.Story(params.ServiceName); err != nil {
		return err
	}

	if err := s.releaseRepo.Store(params.ServiceName, params.ReleaseName, startAt); err != nil {
		return err
	}

	response := apiSpec.RegisterResponse{
		Status: apiSpec.RegisterResponseStatusOk,
	}
	return ctx.JSON(http.StatusOK, response)
}
