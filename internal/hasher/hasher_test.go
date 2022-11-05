package hasher_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/hasher"
)

func TestHashHMAC(t *testing.T) {
	t.Parallel()

	type args struct {
		src string
		key string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test1",
			args: args{
				src: fmt.Sprintf("%s:counter:%d", "test", 10),
				key: "test_super",
			},
			want: "2a38595d8417309951d16599f7202f2ec02c2cdf36612091c1d2e0c2a5d2420f",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := hasher.HashHMAC(tt.args.src, tt.args.key)
			if got != tt.want {
				t.Errorf("HashHMAC() = %v, want %v", got, tt.want)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
