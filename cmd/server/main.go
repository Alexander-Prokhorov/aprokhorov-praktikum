package main

import (
	"aprokhorov-praktikum/cmd/server/handlers"
	"aprokhorov-praktikum/cmd/server/storage"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	yaml "gopkg.in/yaml.v3"
)

func errHandle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type Config struct {
	Address string `yaml:"ADDRESS"`
	Port    string `yaml:"PORT"`
}

func (c *Config) getConfig() *Config {
	yaml_file, err := ioutil.ReadFile("config/config.yaml")
	errHandle(err)
	err = yaml.Unmarshal(yaml_file, c)
	errHandle(err)
	return c
}

func main() {
	var conf Config
	conf.getConfig()

	database := &storage.Mem_storage{}
	database.Init()

	handleUpdateCounter := handlers.HandlerUpdate{
		Metric_type: "counter",
		Storage:     database,
	}
	http.Handle("/update/counter/", handleUpdateCounter)

	handleUpdateGauge := handlers.HandlerUpdate{
		Metric_type: "gauge",
		Storage:     database,
	}
	http.Handle("/update/gauge/", handleUpdateGauge)

	server := &http.Server{
		Addr: conf.Address + ":" + conf.Port,
	}
	log.Fatal(server.ListenAndServe())
}
