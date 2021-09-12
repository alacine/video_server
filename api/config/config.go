package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Configuration 配置项
type Configuration struct {
	StreamAddr string `json:"stream_addr"`
	Scheduler  string `json:"scheduler_addr"`
}

var configuration *Configuration

func init() {
	configPath := filepath.Join("config", "conf.json")
	file, _ := os.Open(configPath)
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("close config file failed %s", err)
		}
	}()
	decoder := json.NewDecoder(file)
	configuration = &Configuration{}
	err := decoder.Decode(configuration)
	log.Printf("streamserver address: %s", configuration.StreamAddr)
	log.Printf("scheduler address: %s", configuration.Scheduler)
	if err != nil {
		panic(err)
	}
}

// GetStreamAddr 获取 streamserver 服务器地址
func GetStreamAddr() string {
	return configuration.StreamAddr
}

// GetScheduler 获取 scheduler 服务器地址
func GetScheduler() string {
	return configuration.Scheduler
}
