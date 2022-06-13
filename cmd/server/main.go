package main

import (
	"log"
	"net/http"

	"aprokhorov-praktikum/cmd/server/config"
	"aprokhorov-praktikum/cmd/server/handlers"
	"aprokhorov-praktikum/cmd/server/storage"

	"github.com/go-chi/chi"
)

func main() {
	// Init Config
	conf := config.NewServerConfig()

	// Init Storage
	database := storage.NewStorageMem()

	// Init chi Router and setup Handlers
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(database))
		r.Route("/value", func(r chi.Router) {
			r.Post("/", handlers.JSONRead(database))
			r.Get("/{metricType}/{metricName}", handlers.Get(database))
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handlers.JSONUpdate(database))
			r.Post("/{metricType}/{metricName}/{metricValue}", handlers.Post(database))
		})
	})

	// Init Server
	server := &http.Server{
		Addr:    conf.Address + ":" + conf.Port,
		Handler: r,
	}
	log.Fatal(server.ListenAndServe())
}
