package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// global context
var registryContext = make(map[string]interface{})

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

type RegisterObjectCommand struct {
	version       uint
	name          string
	remoteAddress string
}

type LookupCommand struct {
	version uint
	name    string
}

func generateKey(name, version string) string {
	return fmt.Sprintf("%s:%d", command.name, command.version)
}

func registerObject(command RegisterObjectCommand) {
	objectKey := generateKey(command.name, command.version)
	registryContext[objectKey] = command
}

func register(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var registerObjectCommand RegisterObjectCommand
	err := decoder.Decode(&registerObjectCommand)
	if err != nil {
		panic(err)
	}
	registerObject(registerObjectCommand)
	log.Println(registerObjectCommand)
}

func lookup(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var lookupCommand LookupCommand
	err := decoder.Decode(&lookupCommand)
	if err != nil {
		panic(err)
	}
	objectKey := generateKey(lookupCommand.name, lookupCommand.version)
	remoteObject = registryContext[objectKey]
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "applicatoin/json")
	json.NewDecoder(w).Encode(remoteObject)
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/register", register)
	http.HandleFunc("/lookup", lookup)
	http.ListenAndServe(":8080", nil)
}
