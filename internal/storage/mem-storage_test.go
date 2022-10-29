package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/storage"
)

func TestNewStorageMem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *storage.MemStorage
	}{
		{
			name: "Creation of MemStorage",
			want: &storage.MemStorage{
				Metrics: storage.Metrics{
					Gauge:   make(map[string]storage.Gauge),
					Counter: make(map[string]storage.Counter),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := storage.NewStorageMem(); !assert.Equal(t, tt.want.Metrics, got.Metrics) {
				t.Errorf("NewStorageMem() = %v, want %v", got.Metrics, tt.want.Metrics)
			}
		})
	}
}

func TestMemStorage_Write(t *testing.T) {
	t.Parallel()

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
				value:      storage.Counter(1),
			},
			wantErr: false,
		},
		{
			name: "MemStorage Gauge Write Test",
			args: args{
				metricName: "NewMetricGauge",
				value:      storage.Gauge(1.1),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ms := storage.NewStorageMem()
			if err := ms.Write(tt.args.metricName, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Read(t *testing.T) {
	t.Parallel()

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
			want:    storage.Counter(1),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ms := storage.NewStorageMem()
			err := ms.Write("Counter1", storage.Counter(1))
			assert.NoError(t, err, "MemCache Write Error")
			got, err := ms.Read(tt.args.valueType, tt.args.metricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("MemStorage.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
