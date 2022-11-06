package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/server/handlers"
	"aprokhorov-praktikum/internal/storage"
)

func TestPost(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reqURL := "/" + strings.Join([]string{"update", tt.args.url.mtype, tt.args.url.name, tt.args.url.value}, "/")

			// Заполним базу тестовыми данными
			tt.args.s = storage.NewStorageMem()
			for _, value := range tt.args.values {
				err := tt.args.s.Write(context.Background(), value.name, value.value)
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
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, tt.want.code, res.StatusCode)
			databaseValue, err := tt.args.s.Read(context.Background(), tt.args.url.mtype, tt.args.url.name)
			if err != nil {
				t.Error("Can't fetch value from database")
			}
			assert.Equal(t, tt.want.value, databaseValue)
		})
	}
}

func ExamplePost() {
	// How to use Post handler
	// Init any storage that compatible with storage.Storage interface{}
	database := storage.NewStorageMem()

	// Init chi-router
	r := chi.NewRouter()

	// Add Get handler endpoint
	r.Post("/", handlers.Post(database))

	// Init and Run Server
	const defaultReadHeaderTimeout = time.Second * 5
	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
	_ = server.ListenAndServe()
}