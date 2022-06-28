package handlers

import (
	"fmt"
	"net/http"

	"aprokhorov-praktikum/cmd/server/storage"
	"aprokhorov-praktikum/internal/hasher"

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

func readHelper(w http.ResponseWriter, s storage.Reader, m *Metrics, key string) error {
	var hashString string

	value, err := s.Read(m.MType, m.ID)
	if err != nil {
		http.Error(w, "404. Not Found", http.StatusNotFound)
		return err
	}
	switch data := value.(type) {
	case storage.Counter:
		respond := int64(data)
		m.Delta = &respond
		hashString = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case storage.Gauge:
		respond := float64(data)
		m.Value = &respond
		hashString = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		http.Error(w, "500. Internal Server Error", http.StatusInternalServerError)
		return err
	}

	if key != "" {
		m.Hash = hasher.HashHMAC(hashString, key)
	}

	return nil
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

		w.Header().Set("Content-Type", "text/html")

		_, err := w.Write([]byte(htmlPage))
		if err != nil {
			panic(err)
		}
		//json.NewEncoder(w).Encode(s.ReadAll())
	}
}
