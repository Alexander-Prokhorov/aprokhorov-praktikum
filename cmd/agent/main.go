package main

import (
	"fmt"
	"log"
	"time"

	"aprokhorov-praktikum/cmd/agent/config"
	"aprokhorov-praktikum/cmd/agent/poller"
	"aprokhorov-praktikum/cmd/agent/sender"
)

func errHandle(text string, err error) {
	if err != nil {
		log.Printf(text, err)
	}
}

func main() {
	// Init Config
	conf := config.Config{}
	conf.InitDefaults()

	// Init Sender
	send := sender.Sender{Server: conf.Server, Port: conf.Port}
	send.Init()

	// Init Poller
	NewMetrics := new(poller.Poller)
	NewMetrics.Init()

	// Poll and Send
	pollInterval, err := time.ParseDuration(conf.PollInterval)
	errHandle("Config parse error: %s", err)

	sendInterval, err := time.ParseDuration(conf.ReportInterval)
	errHandle("Config parse error: %s", err)

	tickerPoll := time.NewTicker(pollInterval)
	tickerSend := time.NewTicker(sendInterval)
	for {
		select {
		case <-tickerPoll.C:
			err = NewMetrics.PollMemStats(conf.MemStatMetrics)
			errHandle("Poller error: %s", err)

			err := NewMetrics.PollRandomMetric()
			errHandle("Poller error: %s", err)

			counter, err := NewMetrics.Storage.Read("counter", "PollCount")
			errHandle("Poller error: %s", err)

			fmt.Println("Poll Count:", counter)

		case <-tickerSend.C:
			fmt.Println("Send Data to Server")

			for metricType, values := range NewMetrics.Storage.ReadAll() {
				for metricName, metricValue := range values {
					go func(mtype string, name string, value string) {
						err := send.SendMetric(mtype, name, value)
						errHandle("Sender error: %s", err)
					}(metricType, metricName, metricValue)
				}

			}
		}
	}
}
