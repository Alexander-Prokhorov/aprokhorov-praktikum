package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"aprokhorov-praktikum/internal/storage"
)

// Handler for POST update metric value by url-encoded input.
func Post(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const (
			base    = 10
			bitSize = 64
		)

		var req Metrics

		req.MType = chi.URLParam(r, "metricType")
		req.ID = chi.URLParam(r, "metricName")
		urlValue := chi.URLParam(r, "metricValue")

		switch req.MType {
		case Counter:
			newValue, err := strconv.ParseInt(urlValue, base, bitSize)
			if err != nil {
				http.Error(w, fmt.Sprintf("400. Can't parse value to int: %s", urlValue), http.StatusBadRequest)

				return
			}

			req.Delta = &newValue
		case Gauge:
			newValue, err := strconv.ParseFloat(urlValue, bitSize)
			if err != nil {
				http.Error(w, fmt.Sprintf("400. Can't parse value to int: %s", urlValue), http.StatusBadRequest)

				return
			}

			req.Value = &newValue
		}

		err := updateHelper(w, s, &req, "")
		if err != nil {
			http.Error(w, fmt.Sprintf("500. Internal Server Error: %s", err), http.StatusBadRequest)

			return
		}

		w.Header().Set("Content-Type", "plain/text")
		http.Error(w, "", http.StatusOK)
	}
}
