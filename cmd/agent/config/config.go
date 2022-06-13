package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Server         string   `yaml:"SERVER"`
	Port           string   `yaml:"PORT"`
	PollInterval   string   `yaml:"POOL_INTERVAL"`
	SendInterval   string   `yaml:"REPORT_INTERVAL"`
	MemStatMetrics []string `yaml:"MEMSTAT_METRICS"`
}

func NewAgentConfig() *Config {
	var c Config

	var envVar struct {
		Addr           string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
		PollInterval   string `env:"POLL_INTERVAL" envDefault:"2s"`
		ReportInterval string `env:"REPORT_INTERVAL" envDefault:"10s"`
	}
	err := env.Parse(&envVar)
	if err != nil {
		log.Fatal(nil)
	}
	varChain := strings.Split(envVar.Addr, ":")

	c.Server = varChain[0]
	c.Port = varChain[1]

	c.PollInterval = envVar.PollInterval
	c.SendInterval = envVar.ReportInterval
	c.MemStatMetrics = sliceMemStat()

	return &c
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
