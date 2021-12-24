package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kiarashazarnia/minimal-rmi/pkg/rmi"
)

// global context
var registryContext = make(map[string]interface{})
var config = rmi.LoadConfig()

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func generateKey(name string, version uint) string {
	return fmt.Sprintf("%s:%d", name, version)
}

func registerObject(command rmi.RegisterObjectCommand) {
	objectKey := generateKey(command.Name, command.Version)
	registryContext[objectKey] = command
}

func register(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var registerObjectCommand rmi.RegisterObjectCommand
	err := decoder.Decode(&registerObjectCommand)
	if err != nil {
		panic(err)
	}
	registerObject(registerObjectCommand)
	log.Println(registerObjectCommand)
}

func lookup(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var lookupCommand rmi.LookupCommand
	err := decoder.Decode(&lookupCommand)
	if err != nil {
		panic(err)
	}
	objectKey := generateKey(lookupCommand.Name, lookupCommand.Version)
	remoteObject := registryContext[objectKey]
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "applicatoin/json")
	json.NewEncoder(w).Encode(remoteObject)
}

func main() {

	http.HandleFunc("/register", register)
	http.HandleFunc("/lookup", lookup)
	http.ListenAndServe(config.RMI_HOST, nil)
}
