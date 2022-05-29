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
	// Init Config
	conf := Config{
		Address: "127.0.0.1",
		Port:    "8080",
	}

	// Init Storage
	database := &storage.MemStorage{}
	database.Init()

	// Init Handler
	handleUpdate := handlers.HandlerUpdate{
		Storage: database,
	}
	http.Handle("/update/", handleUpdate)

	// Init Server
	server := &http.Server{
		Addr: conf.Address + ":" + conf.Port,
	}
	log.Fatal(server.ListenAndServe())
}
