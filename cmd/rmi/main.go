package main

import (
	"fmt"
	"reflect"
	"json"
	"log"
	"net/http"
)

// global context
registryContext := make(map[string]interface{})

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

type RegisterObjectCommand struct {
	version       uint
	remoteAddress string
}


func registerObject(command *RegisterObjectCommand) {
	registryContext[""]
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

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/register", register)

	http.ListenAndServe(":8080", nil)
}
