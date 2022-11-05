package logger_test

import (
	"testing"

	"aprokhorov-praktikum/internal/logger"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		filename string
		level    int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test of Zap-Logger Creation",
			args: args{
				filename: "test.log",
				level:    1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := logger.NewLogger(tt.args.filename, tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
