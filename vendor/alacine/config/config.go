package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	StreamAddr string `json:"stream_addr"`
	Scheduler  string `json:"scheduler_addr"`
}

var configuration *Configuration

func init() {
	file, _ := os.Open("./conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = &Configuration{}
	err := decoder.Decode(configuration)
	log.Printf("streamserver address: %s", configuration.StreamAddr)
	log.Printf("scheduler address: %s", configuration.Scheduler)
	if err != nil {
		panic(err)
	}
}

func GetStreamAddr() string {
	return configuration.StreamAddr
}

func GetScheduler() string {
	return configuration.Scheduler
}
