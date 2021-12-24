package rmi

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	config *Configuration
)

type Configuration struct {
	rmi_host    string
	remote_host string
}

func LoadConfig() *Configuration {

	if config == nil {
		file, _ := os.Open("config.json")
		defer file.Close()
		decoder := json.NewDecoder(file)
		config = new(Configuration)
		err := decoder.Decode(&config)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println(config.rmi_host)
	}

	return config
}
