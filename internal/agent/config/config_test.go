package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/agent/config"
)

func TestNewAgentConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *config.Config
	}{
		{
			name: "Simple Test Creation Default Config",
			want: &config.Config{
				MemStatMetrics: config.SliceMemStat(),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := config.NewAgentConfig(); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewAgentConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SliceMemStat(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := config.SliceMemStat(); !assert.Equal(t, tt.want, got) {
				t.Errorf("sliceMemStat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
	t.Parallel()

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
				MemStatMetrics: config.SliceMemStat(),
				Address:        "",
				PollInterval:   "",
				SendInterval:   "",
				Key:            "",
				Batch:          true,
				LogLevel:       0,
			},
			want: "{\"Address\":\"\",\"PollInterval\":\"\",\"SendInterval\":\"\",\"Key\":\"\",\"CryptoKey\":\"\"}",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := config.Config{
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
