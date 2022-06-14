package main

import (
	"flag"
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

var conf *config.Config

func init() {
	// Init Config
	conf = config.NewServerConfig()

	// Init flags
	flag.StringVar(&conf.Address, "a", "127.0.0.1:8080", "An ip address for server run")
	flag.StringVar(&conf.StoreInterval, "i", "300s", "Interval for storing Data to file")
	flag.StringVar(&conf.StoreFile, "f", "/tmp/devops-metrics-db.json", "File path to store Data")
	flag.BoolVar(&conf.Restore, "r", true, "Restore Metrics from file?")
}

func main() {
	// Init Flags
	flag.Parse()

	// Init Config from Env
	conf.EnvInit()
	fmt.Println(*conf)

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
		Addr:    conf.Address,
		Handler: r,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	// Handle os.Exit
	<-done
}
