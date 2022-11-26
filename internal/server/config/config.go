package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/caarlos0/env/v6"
)

// Structure for Server Config parameters.
type Config struct {
	ConfigFile    string `json:"-" env:"CONFIG"`
	Address       string `json:"address" env:"ADDRESS"`               // envDefault:"127.0.0.1:8080"`
	StoreInterval string `json:"store_interval" env:"STORE_INTERVAL"` // envDefault:"3s"`
	StoreFile     string `json:"store_file" env:"STORE_FILE"`         // envDefault:"/tmp/devops-metrics-db.json"`
	DatabaseDSN   string `json:"database_dsn" env:"DATABASE_DSN"`     // envDefault:"localhost:5432"`
	Restore       bool   `json:"restore" env:"RESTORE"`               // envDefault:"true"`
	Key           string `json:"-" env:"KEY"`                         // envDefault:""`
	CryptoKey     string `json:"crypto_key" env:"CRYPTO_KEY"`         // envDefault:""`
	TrustedSubnet string `json:"trusted_subnet" env:"TRUSTED_SUBNET"` // envDefault:""`
	LogLevel      int    `json:"-" env:"LOG_LEVEL"`
}

// Create new empty Server Config.
func NewServerConfig() *Config {
	return &Config{}
}

// Fill up Server Config from json config File
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
	case c.StoreInterval == "":
		c.StoreInterval = tCfg.StoreInterval
		fallthrough
	case c.StoreFile == "":
		c.StoreFile = tCfg.StoreFile
		fallthrough
	case c.DatabaseDSN == "":
		c.DatabaseDSN = tCfg.DatabaseDSN
		fallthrough
	case !c.Restore:
		c.Restore = tCfg.Restore
		fallthrough
	case c.CryptoKey == "":
		c.CryptoKey = tCfg.CryptoKey
		fallthrough
	case c.TrustedSubnet == "":
		c.CryptoKey = tCfg.TrustedSubnet
	}

	return nil
}

// Init Server Config values from environment variables.
func (c *Config) EnvInit() {
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
}

// Return string representation of Server Config Data.
// For Stringer interface, used in logging.
func (c Config) String() string {
	// cString, err := json.MarshalIndent(c, "", "    ")
	cString, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(cString)
}
