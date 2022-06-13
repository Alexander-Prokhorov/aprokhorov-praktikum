package sender

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSender_SendMetricURL(t *testing.T) {
	type fields struct {
		Server string
		Port   string
		URL    url.URL
		Client http.Client
	}
	type args struct {
		name  string
		mtype string
		value string
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
			fields:  fields{Server: "127.0.0.1", Port: "8080"},
			args:    args{name: "TestMetric", mtype: "counter", value: "1"},
			wantErr: false,
			want:    want{path: "update/counter/TestMetric/1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test Server Initialization
			testServer := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "This is Test HTTPServer")
					}))
			defer testServer.Close()

			// Get Test Server address:port
			serverPort := strings.Replace(testServer.URL, "http://", "", -1)
			params := strings.Split(serverPort, ":")

			// Init Sender Client
			s := NewAgentSender(params[0], params[1])

			if err := s.SendMetricURL(tt.args.mtype, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, s.URL.Path, tt.want.path)
		})
	}
}

func TestSender_SendMetricJSON(t *testing.T) {
	type fields struct {
		Server string
		Port   string
		URL    url.URL
		Client http.Client
	}
	type args struct {
		name  string
		mtype string
		value string
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
			fields:  fields{Server: "127.0.0.1", Port: "8080"},
			args:    args{name: "TestMetric", mtype: "counter", value: "1"},
			wantErr: false,
			want:    want{path: "update"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Test Server Initialization
			testServer := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "This is Test HTTPServer")
					}))
			defer testServer.Close()

			// Get Test Server address:port
			serverPort := strings.Replace(testServer.URL, "http://", "", -1)
			params := strings.Split(serverPort, ":")

			// Init Sender Client
			s := NewAgentSender(params[0], params[1])

			if err := s.SendMetricJSON(tt.args.mtype, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
