package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"aprokhorov-praktikum/internal/storage"
)

// Handler for GET metric value by json-body.
func JSONRead(s storage.Storage, key string) http.HandlerFunc {
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

		err := readHelper(r.Context(), w, s, &jReq, key)
		if err != nil {
			return
		}

		jRes, err := json.Marshal(jReq)
		if err != nil {
			panic(err)
		}

		_, err = w.Write(jRes)
		if err != nil {
			panic(err)
		}
	}
}
