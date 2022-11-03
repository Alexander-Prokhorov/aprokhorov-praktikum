package config

import (
	"encoding/json"
	"log"

	"github.com/caarlos0/env/v6"
)

// Agent Config Data.
type Config struct {
	MemStatMetrics []string `json:"-"`
	Address        string   `env:"ADDRESS"`         // envDefault:"127.0.0.1:8080"`
	PollInterval   string   `env:"POLL_INTERVAL"`   // envDefault:"2s"`
	SendInterval   string   `env:"REPORT_INTERVAL"` // envDefault:"10s"`
	Key            string   `env:"KEY"`             // envDefault:""`
	Batch          bool     `json:"-" env:"-"`
	LogLevel       int      `json:"-" env:"LOG_LEVEL"`
}

// Fill up Agent Config from environment variables.
func (c *Config) EnvInit() {
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
}

// Return string representation of Agent Config Data.
// For Stringer interface, used in logging.
func (c Config) String() string {
	// cString, err := json.MarshalIndent(c, "", "    ")
	cString, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(cString)
}

// Init empty Agent Config.
func NewAgentConfig() *Config {
	return &Config{
		MemStatMetrics: sliceMemStat(),
		Batch:          true,
	}
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
