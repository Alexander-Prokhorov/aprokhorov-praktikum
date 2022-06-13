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
	conf := config.NewAgentConfig()

	// Init Sender
	send := sender.NewAgentSender(conf.Server, conf.Port)

	// Init Poller
	NewMetrics := poller.NewAgentPoller()

	// Poll and Send
	tickerPoll := time.NewTicker(time.Duration(conf.PollInterval) * time.Second)
	tickerSend := time.NewTicker(time.Duration(conf.SendInterval) * time.Second)
	for {
		select {
		case <-tickerPoll.C:
			err := NewMetrics.PollMemStats(conf.MemStatMetrics)
			errHandle("Poller error: %s", err)

			err = NewMetrics.PollRandomMetric()
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
