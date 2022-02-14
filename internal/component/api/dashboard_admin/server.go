package dashboard_admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/goforbroke1006/onix/domain"
	"github.com/goforbroke1006/onix/pkg/log"
)

func NewServer(
	sourceRepo domain.SourceRepository,
	criteriaRepo domain.CriteriaRepository,
	logger log.Logger,
) *server {
	return &server{
		sourceRepo:   sourceRepo,
		criteriaRepo: criteriaRepo,
		logger:       logger,
	}
}

var (
	_ ServerInterface = &server{}
)

type server struct {
	sourceRepo   domain.SourceRepository
	criteriaRepo domain.CriteriaRepository
	logger       log.Logger
}

func (s server) PostCriteria(ctx echo.Context) error {
	requestBody := CreateCriteriaRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&requestBody); err != nil {
		return err
	}
	criteriaID, err := s.criteriaRepo.Create(
		requestBody.ServiceName, requestBody.Title, requestBody.Selector,
		domain.DynamicDirType(requestBody.ExpectedDir),
		domain.MustParsePullPeriodType(string(requestBody.PullPeriod)))
	if err != nil {
		return err
	}

	s.logger.Info("create new criteria")

	resp := CreateResourceResponse{
		NewId:  fmt.Sprintf("%d", criteriaID),
		Status: CreateResourceResponseStatusOk,
	}
	return ctx.JSON(http.StatusOK, resp)
}
