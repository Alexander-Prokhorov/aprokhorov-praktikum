package poller

import (
	"encoding/json"
	"errors"
	"math/rand"
	"runtime"
	"time"

	"aprokhorov-praktikum/cmd/server/storage"
)

type Poller struct {
	Storage storage.Storage
}

func (p *Poller) Init() {
	p.Storage = new(storage.MemStorage)
	p.Storage.Init()
	err := p.Storage.Write("PollCount", storage.Counter(0))
	if err != nil {
		panic("Can't init storage")
	}
}

func (p *Poller) PollRandomMetric() error {
	rand.Seed(time.Now().UnixNano())
	err := p.Storage.Write("RandomValue", storage.Gauge(rand.Float64()))
	if err != nil {
		return err
	}
	return nil
}

func (p *Poller) PollMemStats(lookupMemStat []string) error {

	// Собираем метрики пакетом runtime
	var metricValue runtime.MemStats
	runtime.ReadMemStats(&metricValue)

	// Переводим struct в map через json (костыль?? но проще чем reflect)
	var mapInterface map[string]interface{}
	jsonMemStats, err := json.Marshal(metricValue)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonMemStats, &mapInterface)
	if err != nil {
		return err
	}

	// Выбираем только интересующие нас метрики
	// и записываем их в хранилище агента
	for _, metric := range lookupMemStat {
		targetMetric, ok := mapInterface[metric]
		if !ok {
			continue
		}
		switch data := targetMetric.(type) {
		case int64:
		case float64:
			err := p.Storage.Write(metric, storage.Gauge(data))
			if err != nil {
				return err
			}
		}
	}

	// Увеличим счетчик
	value, err := p.Storage.Read("counter", "PollCount")
	if err != nil {
		return err
	}
	counter, ok := value.(storage.Counter)
	if !ok {
		return errors.New("Can't update counter, it's not a Counter")
	}
	counter++
	err = p.Storage.Write("PollCount", counter)
	if err != nil {
		return err
	}
	return nil
}
