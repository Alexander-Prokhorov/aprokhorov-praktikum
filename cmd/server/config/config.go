package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address string `env:"ADDRESS"`
	Port    string `yaml:"PORT"`
}

func NewServerConfig() *Config {
	var c Config

	var envVar struct {
		Addr string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	}
	err := env.Parse(&envVar)
	if err != nil {
		log.Fatal(nil)
	}
	varChain := strings.Split(envVar.Addr, ":")

	c.Address = varChain[0]
	c.Port = varChain[1]

	return &c
}
