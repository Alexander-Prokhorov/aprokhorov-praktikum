package storage

import (
	"errors"
	"strconv"
	"sync"
)

type MemStorage struct {
	Metrics struct {
		Gauge   map[string]Gauge
		Counter map[string]Counter
	}
	mutex *sync.RWMutex
}

func NewStorageMem() *MemStorage {
	var ms MemStorage
	ms.Metrics.Gauge = make(map[string]Gauge)
	ms.Metrics.Counter = make(map[string]Counter)
	ms.mutex = &sync.RWMutex{}
	return &ms
}

func (ms *MemStorage) Write(metricName string, value interface{}) error {

	switch data := value.(type) {
	case Counter:
		ms.safeCounterWrite(metricName, data)
	case Gauge:
		ms.safeGaugerWrite(metricName, data)
	default:
		err := errors.New("MemFS: Post(): Only [gauge, counter] type are supported")
		return err
	}
	return nil
}

func (ms *MemStorage) Read(valueType string, metricName string) (interface{}, error) {
	switch valueType {
	case "counter":
		return ms.safeCounterRead(metricName)
	case "gauge":
		return ms.safeGaugeRead(metricName)
	default:
		return nil, errors.New("MemFS: Get(): Only [gauge, counter] type are supported")
	}
}

func (ms *MemStorage) ReadAll() map[string]map[string]string {
	ret := make(map[string]map[string]string)
	ret["counter"] = make(map[string]string)
	ret["gauge"] = make(map[string]string)

	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	for metricName, metricValue := range ms.Metrics.Counter {
		ret["counter"][metricName] = strconv.FormatInt(int64(metricValue), 10)
	}

	for metricName, metricValue := range ms.Metrics.Gauge {
		ret["gauge"][metricName] = strconv.FormatFloat(float64(metricValue), 'f', -1, 64)
	}
	return ret
}

func (ms *MemStorage) safeCounterWrite(metricName string, value Counter) {
	ms.mutex.Lock()
	ms.Metrics.Counter[metricName] = value
	ms.mutex.Unlock()
}

func (ms *MemStorage) safeGaugerWrite(metricName string, value Gauge) {
	ms.mutex.Lock()
	ms.Metrics.Gauge[metricName] = value
	ms.mutex.Unlock()
}

func (ms *MemStorage) safeCounterRead(metricName string) (Counter, error) {
	ms.mutex.RLock()
	value, ok := ms.Metrics.Counter[metricName]
	ms.mutex.RUnlock()
	if !ok {
		return Counter(0), errors.New("value not found")
	}
	return value, nil
}

func (ms *MemStorage) safeGaugeRead(metricName string) (Gauge, error) {
	ms.mutex.RLock()
	value, ok := ms.Metrics.Gauge[metricName]
	ms.mutex.RUnlock()
	if !ok {
		return Gauge(0), errors.New("value not found")
	}
	return value, nil
}
