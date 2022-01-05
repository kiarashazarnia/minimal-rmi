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

func registerObject(command rmi.RegisterObjectCommand) {
	objectKey := rmi.GenerateKey(command.Name, command.Version)
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
	var LookupQuery rmi.LookupQuery
	err := decoder.Decode(&LookupQuery)
	if err != nil {
		panic(err)
	}
	objectKey := rmi.GenerateKey(LookupQuery.Name, LookupQuery.Version)
	remoteObject := registryContext[objectKey]
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "applicatoin/json")
	json.NewEncoder(w).Encode(remoteObject)
}

func main() {
	fmt.Println("0")
	http.HandleFunc("/register", register)
	fmt.Println("1")
	http.HandleFunc("/lookup", lookup)
	log.Println("running server on", config.RMI_HOST)
	http.ListenAndServe(config.RMI_HOST, nil)
	fmt.Println("3")
}
