package storage

type Gauge float64
type Counter int64

type Storage interface {
	Init()
	Reader
	Writer
}

type Reader interface {
	Read(valueType string, metric string) (interface{}, error)
	ReadAll() map[string]map[string]string
}

type Writer interface {
	Write(metric string, value interface{}) error
}
