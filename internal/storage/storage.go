package storage

import "context"

type Gauge float64
type Counter int64

type Storage interface {
	Reader
	Writer
	Pinger
	Closer
}

type Reader interface {
	Read(valueType string, metric string) (interface{}, error)
	ReadAll() (map[string]map[string]string, error)
}

type Writer interface {
	Write(metric string, value interface{}) error
}

type Pinger interface {
	Ping(context.Context) error
}
type Closer interface {
	Close()
}
