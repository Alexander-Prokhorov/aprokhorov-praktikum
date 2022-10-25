package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLogger(tt.args.filename, tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
