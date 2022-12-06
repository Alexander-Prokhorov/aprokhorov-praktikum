package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/caarlos0/env/v6"
)

// Agent Config Data.
type Config struct {
	ConfigFile     string   `json:"-" env:"CONFIG"`
	MemStatMetrics []string `json:"-" env:"-"`
	Address        string   `json:"address" env:"ADDRESS"`                 // envDefault:"127.0.0.1:8080"`
	PollInterval   string   `json:"poll_interval" env:"POLL_INTERVAL"`     // envDefault:"2s"`
	SendInterval   string   `json:"report_interval" env:"REPORT_INTERVAL"` // envDefault:"10s"`
	Key            string   `json:"-" env:"KEY"`                           // envDefault:""`
	CryptoKey      string   `json:"crypto_key" env:"CRYPTO_KEY"`           // envDefault:""`
	Batch          bool     `json:"-" env:"-"`
	LogLevel       int      `json:"-" env:"LOG_LEVEL"`
	GRPC           bool     `json:"-" env:"GRPC_CLIENT"`
}

// Fill up Agent Config from json config File
func (c *Config) LoadFromFile() error {
	data, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}

	tCfg := Config{}
	err = json.Unmarshal(data, &tCfg)
	if err != nil {
		return err
	}

	switch {
	case c.Address == "":
		c.Address = tCfg.Address
		fallthrough
	case c.PollInterval == "":
		c.PollInterval = tCfg.PollInterval
		fallthrough
	case c.SendInterval == "":
		c.SendInterval = tCfg.SendInterval
		fallthrough
	case c.CryptoKey == "":
		c.CryptoKey = tCfg.CryptoKey
	}

	return nil
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
