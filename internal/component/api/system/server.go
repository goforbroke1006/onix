package system

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	apiSpec "github.com/goforbroke1006/onix/api/system"
	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

// NewServer creates new server's handlers implementations instance.
func NewServer(
	serviceRepo domain.ServiceRepository,
	releaseRepo domain.ReleaseRepository,
	logger log.Logger,
) *server { // nolint:revive,golint
	return &server{
		serviceRepo: serviceRepo,
		releaseRepo: releaseRepo,
		logger:      logger,
	}
}

var _ apiSpec.ServerInterface = &server{} // nolint:exhaustivestruct

type server struct {
	serviceRepo domain.ServiceRepository
	releaseRepo domain.ReleaseRepository
	logger      log.Logger
}

func (s server) GetHealthz(ctx echo.Context) error {
	err := ctx.NoContent(http.StatusOK)

	return errors.Wrap(err, "write to echo context failed")
}

func (s server) GetRegister(ctx echo.Context, params apiSpec.GetRegisterParams) error {
	startAt := time.Now().UTC()
	if params.StartAt != nil {
		startAt = time.Unix(*params.StartAt, 0).UTC()
	}

	if err := s.serviceRepo.Store(params.ServiceName); err != nil {
		return errors.Wrap(err, "can't store service in repository")
	}

	if err := s.releaseRepo.Store(params.ServiceName, params.ReleaseName, startAt); err != nil {
		return errors.Wrap(err, "can't get release")
	}

	response := apiSpec.RegisterResponse{
		Status: apiSpec.RegisterResponseStatusOk,
	}
	err := ctx.JSON(http.StatusOK, response)

	return errors.Wrap(err, "write to echo context failed")
}
