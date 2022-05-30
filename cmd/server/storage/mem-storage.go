package storage

import (
	"errors"
	"strconv"
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

func (ms *MemStorage) Write(metricName string, value interface{}) error {

	switch data := value.(type) {
	case Counter:
		ms.Metrics.Counter[metricName] = data
	case Gauge:
		ms.Metrics.Gauge[metricName] = data
	default:
		err := errors.New("MemFS: Post(): Only [gauge, counter] type are supported")
		return err
	}
	return nil
}

func (ms *MemStorage) Read(valueType string, metricName string) (interface{}, error) {
	switch valueType {
	case "counter":
		if value, ok := ms.Metrics.Counter[metricName]; ok {
			return value, nil
		} else {
			return nil, errors.New("value not found")
		}
	case "gauge":
		if value, ok := ms.Metrics.Gauge[metricName]; ok {
			return value, nil
		} else {
			return nil, errors.New("value not found")
		}
	default:
		return nil, errors.New("MemFS: Get(): Only [gauge, counter] type are supported")
	}
}

func (ms *MemStorage) ReadAll() map[string]map[string]string {
	ret := make(map[string]map[string]string)
	ret["counter"] = make(map[string]string)
	ret["gauge"] = make(map[string]string)

	for metricName, metricValue := range ms.Metrics.Counter {
		ret["counter"][metricName] = strconv.FormatInt(int64(metricValue), 10)
	}

	for metricName, metricValue := range ms.Metrics.Gauge {
		ret["gauge"][metricName] = strconv.FormatFloat(float64(metricValue), 'f', -1, 64)
	}
	return ret
}
