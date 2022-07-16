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

	"aprokhorov-praktikum/cmd/agent/config"
	"aprokhorov-praktikum/cmd/agent/poller"
	"aprokhorov-praktikum/cmd/agent/sender"
	"aprokhorov-praktikum/internal/logger"

	"go.uber.org/zap"
)

func errHandle(text string, err error, logger *zap.Logger) {
	if err != nil {
		logger.Error(text + err.Error())
	}
}

func main() {
	// Init Config
	conf := config.NewAgentConfig()

	// Init flags
	flag.StringVar(&conf.Address, "a", "127.0.0.1:8080", "An ip address for server run")
	flag.StringVar(&conf.SendInterval, "r", "10s", "Report Interval")
	flag.StringVar(&conf.PollInterval, "p", "2s", "Poll Interval")
	flag.StringVar(&conf.Key, "k", "", "Key for Hash")
	flag.IntVar(&conf.LogLevel, "l", 1, "Log Level, default:Warning")
	flag.Parse()

	// Init Logger
	logger, err := logger.NewLogger("agent.log", conf.LogLevel)
	if err != nil {
		log.Fatal("cannot initialize zap.logger")
	}

	// Init Config from Env
	conf.EnvInit()
	logger.Info(conf.String())

	// Init Sender
	send := sender.NewAgentSender(conf.Address)

	// Init Poller
	NewMetrics := poller.NewAgentPoller()

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
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Init Context and Sync
	ctxMain, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Make Goroutines
	wg.Add(1)
	go func(ctx context.Context, signal <-chan time.Time, sync chan<- struct{}, wgr *sync.WaitGroup, metrics *poller.Poller, metricList []string, log *zap.Logger) {
		for {
			select {
			case <-signal:
				sync <- struct{}{}
				err := metrics.PollMemStats(metricList)
				if err != nil {
					log.Error("Poller MemStat error: " + err.Error())
				}
				err = metrics.PollRandomMetric()
				if err != nil {
					log.Error("Poller MemStat error: " + err.Error())
				}
				counter, err := metrics.Storage.Read("counter", "PollCount")
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
				err := metrics.PollPsUtil()
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
	go func(ctx context.Context, signal <-chan time.Time, wgr *sync.WaitGroup, s *sender.Sender, metrics *poller.Poller, batchStatus bool, key string, log *zap.Logger) {
		for {
			select {
			case <-signal:
				log.Info("Send Data to Server")

				metricsData, err := metrics.Storage.ReadAll()
				if err != nil {
					log.Error("Can't read mertics from storage: " + err.Error())
				}

				// Обновляем либо батчем, либо по одному
				switch batchStatus {
				case true:
					go func() {
						err = s.SendMetricBatch(metricsData, key)
						if err != nil {
							log.Error("Sender Batch: " + err.Error())
						}
					}()
				case false:
					for metricType, values := range metricsData {
						for metricName, metricValue := range values {
							go func(mtype string, name string, value string) {
								err := send.SendMetric(mtype, name, value, key)
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
	logger.Info("Gracefull Close Finished")
}
