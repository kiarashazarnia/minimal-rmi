package rmi

import (
	"encoding/json"
	"log"
	"os"
)

var (
	config *Configuration
)

type Configuration struct {
	RMI_HOST    string
	REMOTE_HOST string
}

func LoadConfig() *Configuration {

	if config == nil {
		file, _ := os.Open("config.json")
		defer file.Close()
		decoder := json.NewDecoder(file)
		config = new(Configuration)
		err := decoder.Decode(&config)
		if err != nil {
			log.Println("error:", err)
		}
		log.Println(config.RMI_HOST)
	}

	return config
}
