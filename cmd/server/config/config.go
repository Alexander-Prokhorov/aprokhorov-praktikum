package config

import (
	"encoding/json"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address       string `env:"ADDRESS"`        // envDefault:"127.0.0.1:8080"`
	StoreInterval string `env:"STORE_INTERVAL"` // envDefault:"3s"`
	StoreFile     string `env:"STORE_FILE"`     // envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool   `env:"RESTORE"`        // envDefault:"true"`
	Key           string `env:"KEY"`            // envDefault:""`
}

func NewServerConfig() *Config {
	return &Config{}
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
