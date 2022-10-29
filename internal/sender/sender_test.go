package sender_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/cmd/agent/sender"
)

func TestSender_SendMetricURL(t *testing.T) {
	t.Parallel()

	type fields struct {
		Address string
		URL     url.URL
		Client  http.Client
	}

	type args struct {
		name  string
		mtype string
		value string
		key   string
	}

	type want struct {
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    want
	}{
		{
			name:    "First Test Sender",
			fields:  fields{Address: "127.0.0.1:8080"},
			args:    args{name: "TestMetric", mtype: "counter", value: "1", key: ""},
			wantErr: false,
			want:    want{path: "update/counter/TestMetric/1"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test Server Initialization
			testServer := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "This is Test HTTPServer")
					}))
			defer testServer.Close()

			// Get Test Server address:port
			params := strings.Split(testServer.URL, "/")

			// Init Sender Client
			s := sender.NewAgentSender(params[2])

			if err := s.SendMetricURL(tt.args.mtype, tt.args.name, tt.args.value, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, s.URL.Path, tt.want.path)
		})
	}
}

func TestSender_SendMetricJSON(t *testing.T) {
	t.Parallel()

	type fields struct {
		Address string
		URL     url.URL
		Client  http.Client
	}

	type args struct {
		name  string
		mtype string
		value string
		key   string
	}

	type want struct {
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    want
	}{
		{
			name:    "JSON Test Sender",
			fields:  fields{Address: "127.0.0.0:8080"},
			args:    args{name: "TestMetric", mtype: "counter", value: "1", key: ""},
			wantErr: false,
			want:    want{path: "update"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test Server Initialization
			testServer := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "This is Test HTTPServer")
					}))
			defer testServer.Close()

			// Get Test Server address:port
			params := strings.Split(testServer.URL, "/")

			// Init Sender Client
			s := sender.NewAgentSender(params[2])

			if err := s.SendMetricJSON(tt.args.mtype, tt.args.name, tt.args.value, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
