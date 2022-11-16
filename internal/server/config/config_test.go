package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/server/config"
)

func TestNewServerConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *config.Config
	}{
		{
			name: "Test creatrion of Config",
			want: &config.Config{
				Address:       "",
				StoreInterval: "",
				StoreFile:     "",
				DatabaseDSN:   "",
				Restore:       false,
				Key:           "",
				LogLevel:      0,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := config.NewServerConfig(); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
	t.Parallel()

	type fields struct {
		Address       string
		StoreInterval string
		StoreFile     string
		DatabaseDSN   string
		Restore       bool
		Key           string
		LogLevel      int
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "String test for Config",
			fields: fields{
				Address:       "",
				StoreInterval: "",
				StoreFile:     "",
				DatabaseDSN:   "",
				Restore:       false,
				Key:           "",
				LogLevel:      0,
			},
			want: "{\"Address\":\"\",\"StoreInterval\":\"\"," +
				"\"StoreFile\":\"\",\"DatabaseDSN\":\"\",\"Restore\":false,\"Key\":\"\",\"CryptoKey\":\"\"}",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := config.Config{
				Address:       tt.fields.Address,
				StoreInterval: tt.fields.StoreInterval,
				StoreFile:     tt.fields.StoreFile,
				DatabaseDSN:   tt.fields.DatabaseDSN,
				Restore:       tt.fields.Restore,
				Key:           tt.fields.Key,
				LogLevel:      tt.fields.LogLevel,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Config.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
