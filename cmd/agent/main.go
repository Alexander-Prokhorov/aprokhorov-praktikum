package main

import (
	"fmt"
	"strconv"
	"time"

	"aprokhorov-praktikum/cmd/agent/config"
	"aprokhorov-praktikum/cmd/agent/poller"
	"aprokhorov-praktikum/cmd/agent/sender"
)

func errHandle(err error) {
	if err != nil {
		fmt.Println(err)
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
	NewMetrics := new(poller.Metrics)

	// Poll and Send
	pollInterval, err := time.ParseDuration(conf.PollInterval)
	errHandle(err)
	sendInterval, err := time.ParseDuration(conf.ReportInterval)
	errHandle(err)
	tickerPoll := time.NewTicker(pollInterval)
	tickerSend := time.NewTicker(sendInterval)
	for {
		select {
		case <-tickerPoll.C:
			err = NewMetrics.PollMemStats(conf.MemStatMetrics)
			if err != nil {
				fmt.Println("Can't fetch metrics")
			}
			NewMetrics.RandomMetric()
			fmt.Println("Poll Count:", NewMetrics.PollCount)
		case <-tickerSend.C:
			fmt.Println("Send Data to Server")
			for name, fValue := range NewMetrics.MemStatMetrics {
				sValue := strconv.FormatFloat(float64(fValue), 'f', -1, 64)
				go func(name string, value string) {
					err := send.SendMetric(name, "guage", value)
					fmt.Println(err)
				}(name, sValue)
			}

			sValue := strconv.FormatFloat(float64(NewMetrics.RandomValue), 'f', -1, 64)
			go func(name string, value string) {
				err := send.SendMetric(name, "guage", value)
				fmt.Println(err)
			}("RandomValue", sValue)

			sValue = strconv.Itoa(int(NewMetrics.PollCount))
			go func(name string, value string) {
				err := send.SendMetric(name, "counter", value)
				fmt.Println(err)
			}("PollCount", sValue)
		}
	}

}
