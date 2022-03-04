package system

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

func NewServer(serviceRepo domain.ServiceRepository, releaseRepo domain.ReleaseRepository, logger log.Logger) *server {
	return &server{
		serviceRepo: serviceRepo,
		releaseRepo: releaseRepo,
		logger:      logger,
	}
}

var (
	_ ServerInterface = &server{}
)

type server struct {
	serviceRepo domain.ServiceRepository
	releaseRepo domain.ReleaseRepository
	logger      log.Logger
}

func (s server) GetHealthz(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (s server) GetRegister(ctx echo.Context, params GetRegisterParams) error {
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

	return ctx.JSON(http.StatusOK, RegisterResponse{Status: RegisterResponseStatusOk})
}
