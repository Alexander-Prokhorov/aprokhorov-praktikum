package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"aprokhorov-praktikum/internal/server/handlers"
	"aprokhorov-praktikum/internal/storage"
)

func TestGet(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Заполним базу тестовыми данными
			tt.args.s = storage.NewStorageMem()
			for _, value := range tt.args.values {
				err := tt.args.s.Write(context.Background(), value.name, value.value)
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
			defer func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, w.Body.String())
		})
	}
}

func ExampleGet() {
	// How to use Get handler
	// Init any storage that compatible with storage.Storage interface{}
	database := storage.NewStorageMem()

	// Init chi-router
	r := chi.NewRouter()

	// Add Get handler endpoint
	r.Post("/", handlers.Get(database))

	// Init and Run Server
	const defaultReadHeaderTimeout = time.Second * 5
	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
	_ = server.ListenAndServe()
}
