package poller

import (
	"encoding/json"
	"errors"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"aprokhorov-praktikum/internal/storage"
)

// Poller Data.
type Poller struct {
	Storage storage.Storage
}

// Init of Poller with new MemStorage.
func NewAgentPoller() *Poller {
	var p Poller
	p.Storage = storage.NewStorageMem()

	err := p.Storage.Write("PollCount", storage.Counter(0))
	if err != nil {
		panic("Can't init storage")
	}

	return &p
}

// Generate and Save random metric.
func (p *Poller) PollRandomMetric() error {
	rand.Seed(time.Now().UnixNano())

	err := p.Storage.Write("RandomValue", storage.Gauge(rand.Float64()))
	if err != nil {
		return err
	}

	return nil
}

// Store Current MemStat Metrics in Poller Storage.
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
			err = p.Storage.Write(metric, storage.Gauge(data))
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
		return errors.New("can't update counter, it's not a counter")
	}
	counter++

	err = p.Storage.Write("PollCount", counter)
	if err != nil {
		return err
	}

	return nil
}

// Store current PsUtil Metrics in Poller storage.
func (p *Poller) PollPsUtil() error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	cpuSlice, err := cpu.Percent(time.Second, false)
	if err != nil {
		return err
	}

	cpuData := cpuSlice[0]

	data := map[string]storage.Gauge{
		"TotalMemory":     storage.Gauge(memory.Total),
		"FreeMemory":      storage.Gauge(memory.Free),
		"CPUutilization1": storage.Gauge(cpuData),
	}

	for name, value := range data {
		err = p.Storage.Write(name, value)
		if err != nil {
			return err
		}
	}

	return nil
}
