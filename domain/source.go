package domain

import "context"

type SourceType string

const (
	SourceTypePrometheus = SourceType("prometheus")
	SourceTypeInfluxDB   = SourceType("influxdb")
)

type Source struct {
	ID      string
	Kind    SourceType
	Address string
}

type SourceRepository interface {
	Create(ctx context.Context, id string, kind SourceType, address string) error
	Get(ctx context.Context, id string) (Source, error)
	GetAll(ctx context.Context) ([]Source, error)
}
