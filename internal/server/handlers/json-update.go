package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"aprokhorov-praktikum/internal/storage"
)

// Handler for POST update single metric value by JSON input in BODY.
func JSONUpdate(s storage.Storage, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != contentJSON {
			errorText := fmt.Sprintf("only application/json supported, get %s", r.Header.Get("Content-Type"))
			http.Error(w, errorText, http.StatusNotImplemented)

			return
		}

		w.Header().Set("Content-Type", contentJSON)

		var jReq Metrics

		if err := json.NewDecoder(r.Body).Decode(&jReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		err := updateHelper(r.Context(), w, s, &jReq, key)
		if err != nil {
			return
		}

		http.Error(w, "", http.StatusOK)
	}
}

// Handler for POST update multiple metrics by JSON input body.
func JSONUpdates(s storage.Storage, key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != contentJSON {
			errorText := fmt.Sprintf("only application/json supported, get %s", r.Header.Get("Content-Type"))
			http.Error(w, errorText, http.StatusNotImplemented)

			return
		}

		w.Header().Set("Content-Type", contentJSON)

		var jReq []Metrics

		if err := json.NewDecoder(r.Body).Decode(&jReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		for i := range jReq {
			err := updateHelper(r.Context(), w, s, &jReq[i], key)
			if err != nil {
				return
			}
		}

		http.Error(w, "", http.StatusOK)
	}
}
