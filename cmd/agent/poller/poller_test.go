package poller

import (
	"fmt"
	"testing"

	"aprokhorov-praktikum/cmd/agent/config"
	"aprokhorov-praktikum/cmd/server/storage"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_PollMemStats(t *testing.T) {
	type fields struct {
		poller Poller
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
				poller: Poller{},
			},
			args: args{
				config: config.Config{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.args.config.InitDefaults()
			tt.fields.poller.Init()

			err := tt.fields.poller.PollMemStats(tt.args.config.MemStatMetrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("Poller.PollMemStats() error = %v, wantErr %v", tt.wantErr, err)
			}

			lenWant := len(tt.args.config.MemStatMetrics)

			gaugeValues := tt.fields.poller.Storage.ReadAll()["gauge"]
			fmt.Print(gaugeValues)
			lenGot := len(gaugeValues)
			assert.Equal(t, lenWant, lenGot)

			pollCount, err := tt.fields.poller.Storage.Read("counter", "PollCount")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, storage.Counter(1), pollCount)
		})
	}
}
