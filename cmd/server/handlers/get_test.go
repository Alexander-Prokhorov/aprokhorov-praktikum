package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"aprokhorov-praktikum/cmd/server/handlers"
	"aprokhorov-praktikum/cmd/server/storage"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type values struct {
		name  string
		value interface{}
	}
	type args struct {
		s      storage.Storage
		url    string
		values []values
	}
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "CounterGet",
			args: args{
				s:   new(storage.MemStorage),
				url: "/value/counter/testCounter",
				values: []values{
					{name: "testCounter", value: storage.Counter(42)},
					{name: "testGauge", value: storage.Gauge(24)},
				},
			},
			want: want{code: 200, response: "42"},
		},
		{
			name: "GaugeGet",
			args: args{
				s:   new(storage.MemStorage),
				url: "/value/gauge/testGauge",
				values: []values{
					{name: "testCounter", value: storage.Counter(42)},
					{name: "testGauge", value: storage.Gauge(24)},
				},
			},
			want: want{code: 200, response: "24"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Заполним базу тестовыми данными
			tt.args.s.Init()
			for _, value := range tt.args.values {
				err := tt.args.s.Write(value.name, value.value)
				if err != nil {
					t.Error(err)
				}
			}

			// Создадим тестовый запрос и рекодер
			r := httptest.NewRequest(http.MethodGet, tt.args.url, nil)
			w := httptest.NewRecorder()

			// Init chi Router and setup Handlers
			cr := chi.NewRouter()
			cr.Get("/value/{metricType}/{metricName}", handlers.Get(tt.args.s))
			cr.ServeHTTP(w, r)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, w.Body.String())
		})
	}
}
