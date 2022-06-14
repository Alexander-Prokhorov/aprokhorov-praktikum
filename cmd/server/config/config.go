package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address       string `env:"ADDRESS"`        // envDefault:"127.0.0.1:8080"`
	StoreInterval string `env:"STORE_INTERVAL"` // envDefault:"3s"`
	StoreFile     string `env:"STORE_FILE"`     // envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool   `env:"RESTORE"`        // envDefault:"true"`
}

func NewServerConfig() *Config {
	var c Config

	err := env.Parse(&c)
	if err != nil {
		log.Fatal(nil)
	}

	return &c
}
