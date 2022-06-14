package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	MemStatMetrics []string
	Server         string
	Port           string
	Address        string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	PollInterval   string `env:"POLL_INTERVAL" envDefault:"2s"`
	SendInterval   string `env:"REPORT_INTERVAL" envDefault:"10s"`
}

func NewAgentConfig() *Config {
	var c Config

	err := env.Parse(&c)
	if err != nil {
		log.Fatal(nil)
	}
	varChain := strings.Split(c.Address, ":")

	c.Server = varChain[0]
	c.Port = varChain[1]

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
