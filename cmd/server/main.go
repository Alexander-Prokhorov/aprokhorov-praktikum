package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"aprokhorov-praktikum/internal/ccrypto"
	"aprokhorov-praktikum/internal/logger"
	"aprokhorov-praktikum/internal/server/config"
	"aprokhorov-praktikum/internal/server/files"
	"aprokhorov-praktikum/internal/server/handlers"
	"aprokhorov-praktikum/internal/storage"
)

// go run -ldflags "-X main.buildVersion=1.1.1 \
// -X 'main.buildDate=$(date +'%Y/%m/%d')' \
// -X 'main.buildCommit=$(git log -1 --pretty=%B | cat)'" \
// main.go
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	// send to stdout buildVars
	if _, err := fmt.Fprintf(
		os.Stdout,
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion,
		buildDate,
		buildCommit,
	); err != nil {
		log.Fatal(err)
	}

	const (
		defaultReadHeaderTimeout = time.Second * 5
	)

	conf := config.NewServerConfig()
	// Init Flags
	flag.StringVar(&conf.ConfigFile, "c", "", "Path to Config File")
	flag.StringVar(&conf.ConfigFile, "config", "", "Path to Config File")
	flag.StringVar(&conf.Address, "a", "127.0.0.1:8080", "An ip address for server run")
	flag.StringVar(&conf.StoreInterval, "i", "300s", "Interval for storing Data to file")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "Path to PostgresSQL (in prefer to File storing)")
	flag.StringVar(&conf.StoreFile, "f", "/tmp/devops-metrics-db.json", "File path to store Data")
	flag.StringVar(&conf.Key, "k", "", "Hash Key")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "", "Path to id_rsa file")
	flag.StringVar(&conf.TrustedSubnet, "t", "", "Trusted Subnet (X-Real-IP check)")
	flag.BoolVar(&conf.Restore, "r", true, "Restore Metrics from file?")
	flag.IntVar(&conf.LogLevel, "l", 1, "Log Level, default:Warning")
	flag.Parse()

	// Init Logger
	logger, err := logger.NewLogger("server.log", conf.LogLevel)
	if err != nil {
		log.Fatal("cannot initialize zap.logger")
	}

	// Init Config from File
	if conf.ConfigFile != "" {
		if err = conf.LoadFromFile(); err != nil {
			logger.Error(fmt.Sprintf("config: cannot load config from file: %s", err.Error()))
		}
	}

	// Init Config from Env
	conf.EnvInit()

	logger.Info(conf.String())

	ctx := context.Background()

	_, acl, err := net.ParseCIDR(conf.TrustedSubnet)
	if err != nil {
		logger.Error(fmt.Sprintf("IP Check is Disabled bacause of error: %s", err.Error()))
	}

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

	// Get Private Key if set
	var privKey *ccrypto.PrivateKey
	if conf.CryptoKey != "" {
		privKey, err = ccrypto.NewPrivateKeyFromFile(conf.CryptoKey)
		if err != nil {
			logger.Error("Failed to load Private Key: " + err.Error())
		}
	}

	// Init chi Router and setup Handlers
	r := chi.NewRouter()

	r.Use(handlers.CheckACL(acl))
	r.Use(handlers.Unpack)
	r.Use(handlers.Pack)
	r.Mount("/debug", middleware.Profiler())
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(database))
		r.Get("/ping", handlers.Ping(database))

		r.Route("/value", func(r chi.Router) {
			r.Use(handlers.Decrypt(privKey))
			r.Post("/", handlers.JSONRead(database, conf.Key))
			r.Get("/{metricType}/{metricName}", handlers.Get(database))
		})

		r.Route("/update", func(r chi.Router) {
			r.Use(handlers.Decrypt(privKey))
			r.Post("/", handlers.JSONUpdate(database, conf.Key))
			r.Post("/{metricType}/{metricName}/{metricValue}", handlers.Post(database))
		})
		r.Route("/updates", func(r chi.Router) {
			r.Use(handlers.Decrypt(privKey))
			r.Post("/", handlers.JSONUpdates(database, conf.Key))
		})
	})

	// Init system calls
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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
