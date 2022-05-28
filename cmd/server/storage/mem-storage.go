package storage

import (
	"errors"
)

type Mem_storage struct {
	Metrics struct {
		Gauge   map[string]Gauge
		Counter map[string]Counter
	}
}

func (ms *Mem_storage) Init() {
	ms.Metrics.Gauge = make(map[string]Gauge)
	ms.Metrics.Counter = make(map[string]Counter)
}

func (ms *Mem_storage) Post(name string, value interface{}) error {

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

func (ms *Mem_storage) Get(value_type string, name string) (interface{}, error) {
	switch value_type {
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
