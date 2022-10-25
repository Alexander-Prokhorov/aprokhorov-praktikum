package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorageMem(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "Creation of MemStorage",
			want: &MemStorage{
				Metrics: Metrics{
					Gauge:   make(map[string]Gauge),
					Counter: make(map[string]Counter),
				},
				mutex: &sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorageMem(); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewStorageMem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_Write(t *testing.T) {

	type args struct {
		metricName string
		value      interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "MemStorage Counter Write Test",
			args: args{
				metricName: "NewMetric",
				value:      Counter(1),
			},
			wantErr: false,
		},
		{
			name: "MemStorage Gauge Write Test",
			args: args{
				metricName: "NewMetricGauge",
				value:      Gauge(1.1),
			},
			wantErr: false,
		},
		{
			name: "MemStorage Fail Write",
			args: args{
				metricName: "NewMetricBool",
				value:      true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewStorageMem()
			if err := ms.Write(tt.args.metricName, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Read(t *testing.T) {

	type args struct {
		valueType  string
		metricName string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Read Counter Test",
			args: args{
				metricName: "Counter1",
				valueType:  "counter",
			},
			want:    Counter(1),
			wantErr: false,
		},
		{
			name: "Read Error Test",
			args: args{
				metricName: "Counter1",
				valueType:  "bool",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewStorageMem()
			err := ms.Write("Counter1", Counter(1))
			assert.NoError(t, err, "MemCache Write Error")
			got, err := ms.Read(tt.args.valueType, tt.args.metricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("MemStorage.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
