package poller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_PollMemStats(t *testing.T) {
	type fields struct {
		MemStatMetrics map[string]gauge
		PollCount      counter
		RandomValue    gauge
	}
	type args struct {
		lookupMemStat []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test for empty Data",
			fields: fields{
				MemStatMetrics: make(map[string]gauge),
				PollCount:      counter(0),
				RandomValue:    gauge(0),
			},
			args: args{
				lookupMemStat: []string{
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				MemStatMetrics: tt.fields.MemStatMetrics,
				PollCount:      tt.fields.PollCount,
				RandomValue:    tt.fields.RandomValue,
			}
			err := m.PollMemStats(tt.args.lookupMemStat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Metrics.PollMemStats() error = %v, wantErr %v", err, tt.wantErr)
			}
			lenWant := len(tt.args.lookupMemStat)
			lenGot := len(m.MemStatMetrics)
			assert.Equal(t, lenGot, lenWant)
		})
	}
}
