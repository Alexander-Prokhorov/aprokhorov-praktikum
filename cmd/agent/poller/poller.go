package poller

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"time"
)

type gauge float64

type counter int64

type Metrics struct {
	MemStatMetrics map[string]gauge
	PollCount      counter
	RandomValue    gauge
}

func (m *Metrics) RandomMetric() {
	rand.Seed(time.Now().UnixNano())
	m.RandomValue = gauge(rand.Float64())
}

func (m *Metrics) PollMemStats(lookupMemStat []string) error {
	var poll_error error
	m.MemStatMetrics = make(map[string]gauge)
	// Собираем метрики пакетом runtime
	var metric_value runtime.MemStats
	runtime.ReadMemStats(&metric_value)
	// Переводим struct в map через json (костыль?? но проще чем reflect)
	var mapInterface map[string]interface{}
	jsonMemStats, err := json.Marshal(metric_value)
	if err != nil {
		poll_error = err
	}
	json.Unmarshal(jsonMemStats, &mapInterface)
	// Выбираем только интересующие нас метрики
	// Сразу конвертруем их в gauge-тип
	for _, metric := range lookupMemStat {
		target_metric := mapInterface[metric]
		switch target_metric.(type) {
		case int64:
			m.MemStatMetrics[metric] = gauge(target_metric.(int64))
		case float64:
			m.MemStatMetrics[metric] = gauge(target_metric.(float64))
		}
	}
	m.PollCount++
	return poll_error
}
