package main

import (
	"flag"
	"fmt"
	"log"
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

	for {
		select {
		case <-tickerPoll.C:
			err := NewMetrics.PollMemStats(conf.MemStatMetrics)
			errHandle("Poller error: %s", err, logger)

			err = NewMetrics.PollRandomMetric()
			errHandle("Poller error: %s", err, logger)

			counter, err := NewMetrics.Storage.Read("counter", "PollCount")
			errHandle("Poller error: %s", err, logger)

			logger.Info(fmt.Sprintf("Poll Count: %v", counter))

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
}
