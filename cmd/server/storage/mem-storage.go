package storage

import (
	"errors"
)

type MemStorage struct {
	Metrics struct {
		Gauge   map[string]Gauge
		Counter map[string]Counter
	}
}

func (ms *MemStorage) Init() {
	ms.Metrics.Gauge = make(map[string]Gauge)
	ms.Metrics.Counter = make(map[string]Counter)
}

func (ms *MemStorage) Post(name string, value interface{}) error {

	switch data := value.(type) {
	case Counter:
		ms.Metrics.Counter[name] = data
	case Gauge:
		ms.Metrics.Gauge[name] = data
	default:
		err := errors.New("MemFS: Post(): Only [gauge, counter] type are supported")
		return err
	}
	return nil
}

func (ms *MemStorage) Get(valueType string, name string) (interface{}, error) {
	switch valueType {
	case "counter":
		if value, ok := ms.Metrics.Counter[name]; ok {
			return value, nil
		} else {
			return Counter(0), nil
		}
	case "gauge":
		if value, ok := ms.Metrics.Gauge[name]; ok {
			return value, nil
		} else {
			return nil, errors.New("value not found")
		}
	default:
		return nil, errors.New("MemFS: Get(): Only [gauge, counter] type are supported")
	}
}
