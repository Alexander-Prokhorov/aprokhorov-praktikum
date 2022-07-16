package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
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

	// Poll and Send
	pollInterval, err := time.ParseDuration(conf.PollInterval)
	errHandle("Config parse error: %s", err, logger)

	sendInterval, err := time.ParseDuration(conf.SendInterval)
	errHandle("Config parse error: %s", err, logger)

	tickerPoll := time.NewTicker(pollInterval)
	tickerSend := time.NewTicker(sendInterval)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context, signal <-chan time.Time, wgr *sync.WaitGroup, metrics *poller.Poller, metricList []string, log *zap.Logger) {
		for {
			select {
			case <-signal:
				err := metrics.PollMemStats(metricList)
				if err != nil {
					log.Error("Poller error: " + err.Error())
				}
				err = metrics.PollRandomMetric()
				if err != nil {
					log.Error("Poller error: " + err.Error())
				}
				counter, err := metrics.Storage.Read("counter", "PollCount")
				if err != nil {
					log.Error("Poller error: " + err.Error())
				}
				log.Info(fmt.Sprintf("Poll Count: %v", counter))
			case <-ctx.Done():
				wgr.Done()
				return
			}
		}
	}(context.Background(), tickerPoll.C, wg, NewMetrics, conf.MemStatMetrics, logger)

	wg.Add(1)
	go func(ctx context.Context, signal <-chan time.Time, s *sender.Sender, metrics *poller.Poller, batchStatus bool, key string, log *zap.Logger) {
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
				wg.Done()
				return
			}
		}

	}(context.Background(), tickerSend.C, send, NewMetrics, conf.Batch, conf.Key, logger)

	/*
		for {
			select {
			case <-tickerPoll.C:

			case <-tickerSend.C:
				logger.Info("Send Data to Server")

				metrics, err := NewMetrics.Storage.ReadAll()
				errHandle("can't read metrics from storage: %s", err, logger)

				// Обновляем либо батчем, либо по одному
				switch conf.Batch {
				case true:
					go func(metric map[string]map[string]string, key string) {
						err = send.SendMetricBatch(metric, key)
						errHandle("Sender Batch error: %s", err, logger)
					}(metrics, conf.Key)
				case false:
					for metricType, values := range metrics {
						for metricName, metricValue := range values {
							go func(mtype string, name string, value string) {
								err := send.SendMetric(mtype, name, value, conf.Key)
								errHandle("Sender error: %s", err, logger)
							}(metricType, metricName, metricValue)
						}
					}
				}

			}
		}
	*/
	wg.Wait()
}
