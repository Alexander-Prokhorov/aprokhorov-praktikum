package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aprokhorov-praktikum/cmd/server/handlers"
	"aprokhorov-praktikum/cmd/server/storage"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	type url struct {
		mtype string
		name  string
		value string
	}
	type values struct {
		name  string
		value interface{}
	}
	type args struct {
		s      storage.Storage
		url    url
		values []values
	}
	type want struct {
		code  int
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "CounterAppend",
			args: args{
				s:   new(storage.MemStorage),
				url: url{mtype: "counter", name: "testCounter", value: "21"},
				values: []values{
					{name: "testCounter", value: storage.Counter(21)},
					{name: "testGauge", value: storage.Gauge(24)},
				},
			},
			want: want{code: 200, value: storage.Counter(42)},
		},
		{
			name: "CounterCreate",
			args: args{
				s:   new(storage.MemStorage),
				url: url{mtype: "counter", name: "newCounter", value: "11"},
				values: []values{
					{name: "testCounter", value: storage.Counter(42)},
					{name: "testGauge", value: storage.Gauge(24)},
				},
			},
			want: want{code: 200, value: storage.Counter(11)},
		},
		{
			name: "GaugeAppend",
			args: args{
				s:   new(storage.MemStorage),
				url: url{mtype: "gauge", name: "testGauge", value: "11.1"},
				values: []values{
					{name: "testCounter", value: storage.Counter(42)},
					{name: "testGauge", value: storage.Gauge(24)},
				},
			},
			want: want{code: 200, value: storage.Gauge(11.1)},
		},
		{
			name: "GaugeCreate",
			args: args{
				s:   new(storage.MemStorage),
				url: url{mtype: "gauge", name: "newGauge", value: "11.1"},
				values: []values{
					{name: "testCounter", value: storage.Counter(42)},
					{name: "testGauge", value: storage.Gauge(24)},
				},
			},
			want: want{code: 200, value: storage.Gauge(11.1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqURL := "/" + strings.Join([]string{"update", tt.args.url.mtype, tt.args.url.name, tt.args.url.value}, "/")

			// Заполним базу тестовыми данными
			tt.args.s = storage.NewStorageMem()
			for _, value := range tt.args.values {
				err := tt.args.s.Write(value.name, value.value)
				if err != nil {
					t.Error(err)
				}
			}

			// Создадим тестовый запрос и рекодер
			r := httptest.NewRequest(http.MethodPost, reqURL, nil)
			w := httptest.NewRecorder()

			// Init chi Router and setup Handlers
			cr := chi.NewRouter()
			cr.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.Post(tt.args.s))
			cr.ServeHTTP(w, r)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			databaseValue, err := tt.args.s.Read(tt.args.url.mtype, tt.args.url.name)
			if err != nil {
				t.Error("Can't fetch value from database")
			}
			assert.Equal(t, tt.want.value, databaseValue)
			//}

		})
	}
}
