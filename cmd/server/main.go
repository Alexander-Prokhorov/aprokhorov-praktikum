package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aprokhorov-praktikum/cmd/server/config"
	"aprokhorov-praktikum/cmd/server/files"
	"aprokhorov-praktikum/cmd/server/handlers"
	"aprokhorov-praktikum/cmd/server/storage"

	"github.com/go-chi/chi"
)

func main() {
	// Init Config
	conf := config.NewServerConfig()

	// Init Storage
	database := storage.NewStorageMem()
	if conf.Restore {
		if err := files.LoadData(conf.StoreFile, database); err != nil {
			log.Fatal(err)
		}
	}

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

	// Init system calls
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Init Saver
	storeInterval, err := time.ParseDuration(conf.StoreInterval)
	if err != nil {
		log.Fatal(err)
	}
	tickerSave := time.NewTicker(storeInterval)
	go func() {
		for {
			<-tickerSave.C
			err := files.SaveData(conf.StoreFile, database)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// defer Save on Exit
	defer func() {
		err := files.SaveData(conf.StoreFile, database)
		if err != nil {
			fmt.Printf("Graceful Shutdown error: %v", err)
		} else {
			fmt.Println("Graceful Shutdown Success!")
		}
	}()
	defer fmt.Println("\nGraceful Shutdown Started!")

	// Init Server
	server := &http.Server{
		Addr:    conf.Server + ":" + conf.Port,
		Handler: r,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	// Handle os.Exit
	<-done
}
