package handlers

import (
	"fmt"
	"net/http"

	"aprokhorov-praktikum/internal/storage"

	"github.com/go-chi/chi"
)

func Get(s storage.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		var req Metrics

		req.MType = chi.URLParam(r, "metricType")
		req.ID = chi.URLParam(r, "metricName")

		err := readHelper(w, s, &req, "")
		if err != nil {
			return
		}

		var respond interface{}
		switch req.MType {
		case "counter":
			respond = *req.Delta
		case "gauge":
			respond = *req.Value
		}

		_, err = w.Write([]byte(fmt.Sprintf("%v", respond)))
		if err != nil {
			panic(err)
		}
	}
}

func GetAll(s storage.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decorator := func(text string, htmlTag string) string {
			return "<" + htmlTag + ">" + text + "</" + htmlTag + ">"
		}

		var htmlPage string
		htmlPage += decorator("All Metrics", "h1")
		metrics, err := s.ReadAll()
		if err != nil {
			http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
			return
		}
		for metricType, metrics := range metrics {
			htmlPage += decorator(metricType, "h2")
			for metricName, MetricValue := range metrics {
				htmlPage += decorator(metricName+" = "+MetricValue, "div")
			}
		}

		w.Header().Set("Content-Type", "text/html")

		_, err = w.Write([]byte(htmlPage))
		if err != nil {
			panic(err)
		}
		//json.NewEncoder(w).Encode(s.ReadAll())
	}
}
