package storage

type Gauge float64
type Counter int64

type Storage interface {
	Init()
	Post(metric string, value interface{}) error
	Get(value_type string, metric string) (interface{}, error)
}
