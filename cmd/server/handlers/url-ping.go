package handlers

import (
	"context"
	"fmt"
	"net/http"

	"aprokhorov-praktikum/internal/storage"
)

// Handler for GET service health check
func Ping(s storage.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		err := s.Ping(context.Background())
		if err != nil {
			http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte("success"))
		if err != nil {
			http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
			return
		}
	}
}
