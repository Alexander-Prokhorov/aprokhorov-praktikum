package config

import (
	"encoding/json"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	MemStatMetrics []string `json:"-"`
	Address        string   `env:"ADDRESS"`         //envDefault:"127.0.0.1:8080"`
	PollInterval   string   `env:"POLL_INTERVAL"`   //envDefault:"2s"`
	SendInterval   string   `env:"REPORT_INTERVAL"` //envDefault:"10s"`
}

func (c *Config) EnvInit() {
	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Config) String() string {
	cString, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return ""
	}
	return string(cString)
}

func NewAgentConfig() *Config {
	return &Config{MemStatMetrics: sliceMemStat()}
}

func sliceMemStat() []string {
	return []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}
}
