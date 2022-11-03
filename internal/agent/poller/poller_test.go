package poller_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/agent/config"
	"aprokhorov-praktikum/internal/agent/poller"
	"aprokhorov-praktikum/internal/storage"
)

func TestMetrics_PollMemStats(t *testing.T) {
	t.Parallel()

	type fields struct {
		poller poller.Poller
	}

	type args struct {
		config config.Config
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
				poller: poller.Poller{},
			},
			args: args{
				config: config.Config{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			tt.args.config = *config.NewAgentConfig()
			tt.fields.poller = *poller.NewAgentPoller(ctx)

			err := tt.fields.poller.PollMemStats(ctx, tt.args.config.MemStatMetrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("Poller.PollMemStats() error = %v, wantErr %v", tt.wantErr, err)
			}

			lenWant := len(tt.args.config.MemStatMetrics)

			data, err := tt.fields.poller.Storage.ReadAll(ctx)
			assert.NoError(t, err)
			gaugeValues := data["gauge"]
			lenGot := len(gaugeValues)
			assert.Equal(t, lenWant, lenGot)

			pollCount, err := tt.fields.poller.Storage.Read(ctx, "counter", "PollCount")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, storage.Counter(1), pollCount)
		})
	}
}
