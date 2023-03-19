package impl

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/goforbroke1006/onix/internal/component/api/external/v1/spec"
)

func (h handlersImpl) PostRelease(ctx echo.Context, params spec.PostReleaseParams) error {
	startAt := time.Now().UTC()
	if params.StartAt != nil {
		startAt = time.Unix(*params.StartAt, 0).UTC()
	}
	storeErr := h.releaseRepo.Store(params.Service, params.Release, startAt)
	if storeErr != nil {
		zap.L().Error("store release fail", zap.Error(storeErr))
		return errors.Wrap(storeErr, "can't get releases list")
	}

	return ctx.String(http.StatusCreated, "OK")
}

func (h handlersImpl) GetRelease(ctx echo.Context, params spec.GetReleaseParams) error {
	ranges, err := h.releaseSvc.GetAll(params.Service)
	if err != nil {
		zap.L().Error("load releases list fail", zap.Error(err))
		return errors.Wrap(err, "can't get releases list")
	}

	response := make([]spec.Release, 0, len(ranges))
	for _, r := range ranges {
		response = append(response, spec.Release{
			Id:    r.ID,
			Title: r.Tag,
			From:  r.StartAt.Unix(),
			Till:  r.StopAt.Unix(),
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
