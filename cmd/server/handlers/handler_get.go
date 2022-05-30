package handlers

import (
	"aprokhorov-praktikum/cmd/server/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func Get(s storage.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		value, err := s.Read(metricType, metricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404. Not Found"))
			return
		}
		var respond string
		switch data := value.(type) {
		case storage.Counter:
			respond = strconv.FormatInt(int64(data), 10)
		case storage.Gauge:
			respond = strconv.FormatFloat(float64(data), 'f', -1, 64)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500. Internal Server Error"))
			return
		}
		w.Write([]byte(respond))
	}
}

func GetAll(s storage.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decorator := func(text string, htmlTag string) string {
			return "<" + htmlTag + ">" + text + "</" + htmlTag + ">"
		}

		var htmlPage string
		htmlPage += decorator("All Metrics", "h1")
		for metricType, metrics := range s.ReadAll() {
			htmlPage += decorator(metricType, "h2")
			for metricName, MetricValue := range metrics {
				htmlPage += decorator(metricName+" = "+MetricValue, "div")
			}
		}
		w.Write([]byte(htmlPage))
		//json.NewEncoder(w).Encode(s.ReadAll())
	}
}
