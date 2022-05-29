package main

import (
	"aprokhorov-praktikum/cmd/agent/poller"
	"aprokhorov-praktikum/cmd/agent/sender"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v3"
)

func errHandle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type Config struct {
	Server         string   `yaml:"SERVER"`
	Port           string   `yaml:"PORT"`
	PollInterval   string   `yaml:"POOL_INTERVAL"`
	ReportInterval string   `yaml:"REPORT_INTERVAL"`
	MemStatMetrics []string `yaml:"MEMSTAT_METRICS"`
}

func (c *Config) getConfig() *Config {
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	errHandle(err)
	err = yaml.Unmarshal(yamlFile, c)
	errHandle(err)
	return c
}

func main() {
	fmt.Println(os.Environ())
	// Init Config
	var conf Config
	conf.getConfig()

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
			NewMetrics.PollMemStats(conf.MemStatMetrics)
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
