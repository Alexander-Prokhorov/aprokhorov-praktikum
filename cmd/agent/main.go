package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"aprokhorov-praktikum/internal/agent/config"
	"aprokhorov-praktikum/internal/agent/poller"
	"aprokhorov-praktikum/internal/agent/sender"
	"aprokhorov-praktikum/internal/ccrypto"
	"aprokhorov-praktikum/internal/logger"
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

func errHandle(text string, err error, logger *zap.Logger) {
	if err != nil {
		logger.Error(text + err.Error())
	}
}

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

	// Init Config
	conf := config.NewAgentConfig()

	// Init flags
	flag.StringVar(&conf.ConfigFile, "c", "", "Path to Config File")
	flag.StringVar(&conf.ConfigFile, "config", "", "Path to Config File")
	flag.StringVar(&conf.Address, "a", "127.0.0.1:8080", "An ip address for server run")
	flag.StringVar(&conf.SendInterval, "r", "10s", "Report Interval")
	flag.StringVar(&conf.PollInterval, "p", "2s", "Poll Interval")
	flag.StringVar(&conf.Key, "k", "", "Key for Hash")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "", "Path to id_rsa.pub file")
	flag.IntVar(&conf.LogLevel, "l", 1, "Log Level, default:Warning")
	flag.Parse()

	// Init Logger
	logger, err := logger.NewLogger("agent.log", conf.LogLevel)
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

	// Init Sender
	send := sender.NewAgentSender(conf.Address)

	// Init Poller
	NewMetrics := poller.NewAgentPoller(ctx)

	// Get Public Key if set
	var pubKey *ccrypto.PublicKey
	if conf.CryptoKey != "" {
		pubKey, err = ccrypto.NewPublicKeyFromFile(conf.CryptoKey)
		if err != nil {
			logger.Fatal("Failed to load Public Key: " + err.Error())
		}
	}

	// Poll and Send tickers
	pollInterval, err := time.ParseDuration(conf.PollInterval)
	errHandle("Config parse error: %s", err, logger)

	sendInterval, err := time.ParseDuration(conf.SendInterval)
	errHandle("Config parse error: %s", err, logger)

	tickerPoll := time.NewTicker(pollInterval)
	tickerSend := time.NewTicker(sendInterval)
	syncChan := make(chan struct{})

	// Init system calls
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Init Context and Sync
	ctxMain, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Make Goroutines
	wg.Add(1)

	go func(
		ctx context.Context,
		signal <-chan time.Time,
		sync chan<- struct{},
		wgr *sync.WaitGroup,
		metrics *poller.Poller,
		metricList []string,
		log *zap.Logger,
	) {
		for {
			select {
			case <-signal:
				sync <- struct{}{}

				err := metrics.PollMemStats(ctx, metricList)
				if err != nil {
					log.Error("Poller MemStat error: " + err.Error())
				}

				err = metrics.PollRandomMetric(ctx)
				if err != nil {
					log.Error("Poller MemStat error: " + err.Error())
				}

				counter, err := metrics.Storage.Read(ctx, "counter", "PollCount")
				if err != nil {
					log.Error("Poller MemStat error: " + err.Error())
				}

				log.Info(fmt.Sprintf("Poll MemStat Count: %v", counter))
			case <-ctx.Done():
				log.Info("Close Poller MemStat Goroutine")
				wgr.Done()

				return
			}
		}
	}(ctxMain, tickerPoll.C, syncChan, wg, NewMetrics, conf.MemStatMetrics, logger)

	wg.Add(1)

	go func(ctx context.Context, signal <-chan struct{}, wgr *sync.WaitGroup, metrics *poller.Poller, log *zap.Logger) {
		for {
			select {
			case <-signal:
				err := metrics.PollPsUtil(ctx)
				if err != nil {
					log.Error("Poller PSUtil error: " + err.Error())
				}

				log.Info("Poll PSUtil Done")
			case <-ctx.Done():
				log.Info("Close Poller PSUtil Goroutine")
				wgr.Done()

				return
			}
		}
	}(ctxMain, syncChan, wg, NewMetrics, logger)

	wg.Add(1)

	go func(
		ctx context.Context,
		signal <-chan time.Time,
		wgr *sync.WaitGroup,
		s *sender.Sender,
		metrics *poller.Poller,
		batchStatus bool,
		key string,
		log *zap.Logger,
	) {
		for {
			select {
			case <-signal:
				log.Info("Send Data to Server")

				metricsData, err := metrics.Storage.ReadAll(ctx)
				if err != nil {
					log.Error("Can't read mertics from storage: " + err.Error())
				}

				// Обновляем либо батчем, либо по одному
				switch batchStatus {
				case true:
					go func() {
						err = s.SendMetricBatch(metricsData, key, pubKey)
						if err != nil {
							log.Error("Sender Batch: " + err.Error())
						}
					}()
				case false:
					for metricType, values := range metricsData {
						for metricName, metricValue := range values {
							go func(mtype string, name string, value string) {
								err := send.SendMetricSingle(mtype, name, value, key, pubKey)
								if err != nil {
									log.Error("Sender Simple: " + err.Error())
								}
							}(metricType, metricName, metricValue)
						}
					}
				}
			case <-ctx.Done():
				log.Info("Close Sender Goroutine")
				wgr.Done()

				return
			}
		}
	}(ctxMain, tickerSend.C, wg, send, NewMetrics, conf.Batch, conf.Key, logger)

	// Handle system calls
	<-done
	cancel()
	wg.Wait()
	logger.Info("Graceful Close Finished")
}
