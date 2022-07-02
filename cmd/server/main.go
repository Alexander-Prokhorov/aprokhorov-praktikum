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
	"aprokhorov-praktikum/internal/storage"

	"github.com/go-chi/chi"
)

func main() {
	conf := config.NewServerConfig()
	// Init Flags
	flag.StringVar(&conf.Address, "a", "127.0.0.1:8080", "An ip address for server run")
	flag.StringVar(&conf.StoreInterval, "i", "300s", "Interval for storing Data to file")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "Path to PostgresSQL (in prefer to File storing)")
	flag.StringVar(&conf.StoreFile, "f", "/tmp/devops-metrics-db.json", "File path to store Data")
	flag.StringVar(&conf.Key, "k", "", "Hash Key")
	flag.BoolVar(&conf.Restore, "r", true, "Restore Metrics from file?")
	flag.Parse()

	// Init Config from Env
	conf.EnvInit()
	fmt.Println(*conf)

	// Init Storage
	var database storage.Storage
	if conf.DatabaseDSN == "" {
		database = storage.NewStorageMem()
		if conf.Restore {
			if err := files.LoadData(conf.StoreFile, database); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		var err error
		database, err = storage.NewDatabaseConnect(conf.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer database.Close()

	// Init chi Router and setup Handlers
	r := chi.NewRouter()

	r.Use(handlers.Unpack)
	r.Use(handlers.Pack)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(database))
		r.Route("/ping", func(r chi.Router) {
			r.Get("/", handlers.Ping(database))
		})
		r.Route("/value", func(r chi.Router) {
			r.Post("/", handlers.JSONRead(database, conf.Key))
			r.Get("/{metricType}/{metricName}", handlers.Get(database))
		})
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handlers.JSONUpdate(database, conf.Key))
			r.Post("/{metricType}/{metricName}/{metricValue}", handlers.Post(database))
		})
	})

	// Init system calls
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if conf.DatabaseDSN == "" {
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
	}
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
