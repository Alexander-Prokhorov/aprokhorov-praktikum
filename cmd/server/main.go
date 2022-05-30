package main

import (
	"aprokhorov-praktikum/cmd/server/handlers"
	"aprokhorov-praktikum/cmd/server/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Config struct {
	Address string `yaml:"ADDRESS"`
	Port    string `yaml:"PORT"`
}

func main() {
	// Init Config
	conf := Config{
		Address: "127.0.0.1",
		Port:    "8080",
	}

	// Init Storage
	database := &storage.MemStorage{}
	database.Init()

	// Init chi Router and setup Handlers
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(database))
		r.Route("/value/{metricType}", func(r chi.Router) {
			r.Get("/{metricName}", handlers.Get(database))
		})
		r.Route("/update/{metricType}", func(r chi.Router) {
			r.Post("/{metricName}/{metricValue}", handlers.Post(database))
		})
	})

	// Init Server
	server := &http.Server{
		Addr:    conf.Address + ":" + conf.Port,
		Handler: r,
	}
	log.Fatal(server.ListenAndServe())
}
