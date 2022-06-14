package main

import (
	"flag"
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

var conf *config.Config

func init() {
	// Init Config
	conf = config.NewAgentConfig()

	// Init flags
	conf.Address = *flag.String("a", "127.0.0.0:8080", "An ip address for server run")
	conf.SendInterval = *flag.String("r", "10s", "Report Interval")
	conf.PollInterval = *flag.String("p", "2s", "Poll Interval")

	// Re-init for Env vars
	conf = config.NewAgentConfig()

}

func main() {
	//Init Flags
	flag.Parse()

	// Init Sender
	send := sender.NewAgentSender(conf.Address)

	// Init Poller
	NewMetrics := poller.NewAgentPoller()

	// Poll and Send
	pollInterval, err := time.ParseDuration(conf.PollInterval)
	errHandle("Config parse error: %s", err)

	sendInterval, err := time.ParseDuration(conf.SendInterval)
	errHandle("Config parse error: %s", err)

	tickerPoll := time.NewTicker(pollInterval)
	tickerSend := time.NewTicker(sendInterval)

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
