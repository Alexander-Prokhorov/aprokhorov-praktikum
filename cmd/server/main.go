package main

import (
	"aprokhorov-praktikum/cmd/server/handlers"
	"aprokhorov-praktikum/cmd/server/storage"
	"fmt"
	"log"
	"net/http"
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

func main() {
	conf := Config{
		Address: "127.0.0.1",
		Port:    "8080",
	}

	database := &storage.MemStorage{}
	database.Init()

	handleUpdateCounter := handlers.HandlerUpdate{
		MetricType: "counter",
		Storage:    database,
	}
	http.Handle("/update/counter/", handleUpdateCounter)

	handleUpdateGauge := handlers.HandlerUpdate{
		MetricType: "gauge",
		Storage:    database,
	}
	http.Handle("/update/gauge/", handleUpdateGauge)

	server := &http.Server{
		Addr: conf.Address + ":" + conf.Port,
	}
	log.Fatal(server.ListenAndServe())
}
