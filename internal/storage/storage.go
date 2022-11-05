package storage

import "context"

type (
	Gauge   float64
	Counter int64
)

type Storage interface {
	Reader
	Writer
	Pinger
	Closer
}

type Reader interface {
	Read(ctx context.Context, valueType string, metric string) (interface{}, error)
	ReadAll(ctx context.Context) (map[string]map[string]string, error)
}

type Writer interface {
	Write(ctx context.Context, metric string, value interface{}) error
}

type Pinger interface {
	Ping(context.Context) error
}
type Closer interface {
	Close()
}
