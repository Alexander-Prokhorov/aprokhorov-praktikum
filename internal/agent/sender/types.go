package sender

import (
	"context"

	"aprokhorov-praktikum/internal/ccrypto"
)

// Metric Data.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type Sender interface {
	SendMetricSingle(context.Context, string, string, string, string, *ccrypto.PublicKey) error
	SendMetricBatch(context.Context, map[string]map[string]string, string, *ccrypto.PublicKey) error
}
