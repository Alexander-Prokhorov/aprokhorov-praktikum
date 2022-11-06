package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"aprokhorov-praktikum/internal/logger"
	"aprokhorov-praktikum/internal/server/config"
	"aprokhorov-praktikum/internal/server/files"
	"aprokhorov-praktikum/internal/server/handlers"
	"aprokhorov-praktikum/internal/storage"
)

func main() {
	const (
		defaultReadHeaderTimeout = time.Second * 5
	)

	conf := config.NewServerConfig()
	// Init Flags
	flag.StringVar(&conf.Address, "a", "127.0.0.1:8080", "An ip address for server run")
	flag.StringVar(&conf.StoreInterval, "i", "300s", "Interval for storing Data to file")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "Path to PostgresSQL (in prefer to File storing)")
	flag.StringVar(&conf.StoreFile, "f", "/tmp/devops-metrics-db.json", "File path to store Data")
	flag.StringVar(&conf.Key, "k", "", "Hash Key")
	flag.BoolVar(&conf.Restore, "r", true, "Restore Metrics from file?")
	flag.IntVar(&conf.LogLevel, "l", 1, "Log Level, default:Warning")
	flag.Parse()

	// Init Config from Env
	conf.EnvInit()

	ctx := context.Background()

	// Init Logger
	logger, err := logger.NewLogger("server.log", conf.LogLevel)
	if err != nil {
		log.Fatal("cannot initialize zap.logger")
	}

	logger.Info(conf.String())

	// Init Storage
	var database storage.Storage

	switch conf.DatabaseDSN {
	case "":
		database = storage.NewStorageMem()
		if conf.Restore {
			err = files.LoadData(conf.StoreFile, database)
			if err != nil {
				logger.Fatal(fmt.Sprintf("can't load data from file: %s", err.Error()))
			}
		}
	default:
		database, err = storage.NewDatabaseConnect(ctx, conf.DatabaseDSN)
		if err != nil {
			logger.Fatal(fmt.Sprintf("can't connect to database: %s", err.Error()))
		}
	}

	defer func() {
		err = database.Close()
	}()

	// Init chi Router and setup Handlers
	r := chi.NewRouter()

	r.Use(handlers.Unpack)
	r.Use(handlers.Pack)
	r.Mount("/debug", middleware.Profiler())
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(database))
		r.Get("/ping", handlers.Ping(database))

		r.Route("/value", func(r chi.Router) {
			r.Post("/", handlers.JSONRead(database, conf.Key))
			r.Get("/{metricType}/{metricName}", handlers.Get(database))
		})

		r.Route("/update", func(r chi.Router) {
			r.Post("/", handlers.JSONUpdate(database, conf.Key))
			r.Post("/{metricType}/{metricName}/{metricValue}", handlers.Post(database))
		})

		r.Post("/updates/", handlers.JSONUpdates(database, conf.Key))
	})

	// Init system calls
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if conf.DatabaseDSN == "" {
		// Init Saver
		storeInterval, err := time.ParseDuration(conf.StoreInterval)
		if err != nil {
			logger.Fatal(fmt.Sprintf("can't parse store inverval: %s", err.Error()))
		}

		tickerSave := time.NewTicker(storeInterval)

		go func() {
			for {
				<-tickerSave.C

				_ = files.SaveData(conf.StoreFile, database)
			}
		}()

		// defer Save on Exit
		defer func() {
			err := files.SaveData(conf.StoreFile, database)
			if err != nil {
				logger.Error(fmt.Sprintf("Graceful Shutdown error: %s", err.Error()))
			} else {
				logger.Info("Graceful Shutdown Success!")
			}
		}()
		defer logger.Info("Graceful Shutdown Started!")
	}
	// Init Server
	server := &http.Server{
		Addr:              conf.Address,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Handler:           r,
	}

	go func() {
		logger.Fatal(server.ListenAndServe().Error())
	}()

	// Handle os.Exit
	<-done
}
