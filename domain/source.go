package domain

type SourceType string

const (
	SourceTypePrometheus = SourceType("prometheus")
	SourceTypeInfluxDB   = SourceType("influxdb")
)

type Source struct {
	ID      int64
	Title   string
	Type    SourceType
	Address string
}

type SourceRepository interface {
	Create(title string, kind SourceType, address string) (int64, error)
	Get(identifier int64) (*Source, error)
	GetAll() ([]Source, error)
}
