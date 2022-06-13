package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"aprokhorov-praktikum/cmd/server/storage"
)

func JsonUpdate(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			error_text := fmt.Sprintf("only application/json supported, get %s", r.Header.Get("Content-Type"))
			http.Error(w, error_text, http.StatusNotImplemented)
			return
		}

		var jReq Metrics

		if err := json.NewDecoder(r.Body).Decode(&jReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := updateHelper(w, s, &jReq)
		if err != nil {
			return
		}

	}
}
