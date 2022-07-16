package logger

import (
	"go.uber.org/zap"
)

func NewLogger(filename string, level int) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if level == int(zap.DebugLevel) {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.OutputPaths = []string{
		filename,
	}

	return cfg.Build()
}
