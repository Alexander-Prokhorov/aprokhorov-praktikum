package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAgentConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Simple Test Creation Default Config",
			want: &Config{
				MemStatMetrics: sliceMemStat(),
				Address:        "",
				PollInterval:   "",
				SendInterval:   "",
				Key:            "",
				Batch:          true,
				LogLevel:       0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgentConfig(); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewAgentConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sliceMemStat(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "Simple Test for wanted Metrics",
			want: []string{
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceMemStat(); !assert.Equal(t, tt.want, got) {
				t.Errorf("sliceMemStat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
	type fields struct {
		MemStatMetrics []string
		Address        string
		PollInterval   string
		SendInterval   string
		Key            string
		Batch          bool
		LogLevel       int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "String test for Config",
			fields: fields{
				MemStatMetrics: sliceMemStat(),
				Address:        "",
				PollInterval:   "",
				SendInterval:   "",
				Key:            "",
				Batch:          true,
				LogLevel:       0,
			},
			want: "{\"Address\":\"\",\"PollInterval\":\"\",\"SendInterval\":\"\",\"Key\":\"\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				MemStatMetrics: tt.fields.MemStatMetrics,
				Address:        tt.fields.Address,
				PollInterval:   tt.fields.PollInterval,
				SendInterval:   tt.fields.SendInterval,
				Key:            tt.fields.Key,
				Batch:          tt.fields.Batch,
				LogLevel:       tt.fields.LogLevel,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Config.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
