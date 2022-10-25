package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Test creatrion of Config",
			want: &Config{
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
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServerConfig(); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewServerConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
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
			want: "{\"Address\":\"\",\"StoreInterval\":\"\",\"StoreFile\":\"\",\"DatabaseDSN\":\"\",\"Restore\":false,\"Key\":\"\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
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
