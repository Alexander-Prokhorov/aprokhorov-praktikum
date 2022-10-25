package config

import (
	"encoding/json"
	"log"

	"github.com/caarlos0/env/v6"
)

// Structure for Server Config parameters.
type Config struct {
	Address       string `env:"ADDRESS"`        // envDefault:"127.0.0.1:8080"`
	StoreInterval string `env:"STORE_INTERVAL"` // envDefault:"3s"`
	StoreFile     string `env:"STORE_FILE"`     // envDefault:"/tmp/devops-metrics-db.json"`
	DatabaseDSN   string `env:"DATABASE_DSN"`   // envDefault:"localhost:5432"`
	Restore       bool   `env:"RESTORE"`        // envDefault:"true"`
	Key           string `env:"KEY"`            // envDefault:""`
	LogLevel      int    `json:"-" env:"LOG_LEVEL"`
}

// Create new empty Server Config
func NewServerConfig() *Config {
	return &Config{}
}

// Init Server Config values from environment variables.
func (c *Config) EnvInit() {
	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}
}

// Return string representation of Server Config Data.
// For Stringer interface, used in logging.
func (c Config) String() string {
	//cString, err := json.MarshalIndent(c, "", "    ")
	cString, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(cString)
}
