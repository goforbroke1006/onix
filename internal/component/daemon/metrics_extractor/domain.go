package metrics_extractor

import (
	"github.com/goforbroke1006/onix/domain"
)

type ReleaseRepository interface {
	GetLast(serviceName string) (*domain.Release, error)
}
