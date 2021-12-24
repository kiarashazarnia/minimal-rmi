package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/kiarashazarnia/minimal-rmi/pkg/rmi"
)

func remote(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var methodCall rmi.MethodCall
	err := decoder.Decode(&methodCall)
	if err != nil {
		panic(err)
	}
	handleMethodCall(methodCall)

	log.Println(methodCall)
}

var config = rmi.LoadConfig()

func main() {

	var hello Hello = new(HelloRemoteObject)

	register(hello, 1)

	http.HandleFunc("/remote", remote)
	http.ListenAndServe(config.REMOTE_HOST, nil)
}

func register(object interface{}, version uint) bool {

	data := rmi.RegisterObjectCommand{
		Version:       version,
		Name:          reflect.TypeOf(object).Name(),
		RemoteAddress: config.REMOTE_HOST,
	}
	jsonData, _ := json.Marshal(data)
	response, _ := http.Post(config.RMI_HOST, "application/json", bytes.NewBuffer(jsonData))
	return response.StatusCode == http.StatusOK
}

type HelloRemoteObject struct {
	helloSentence string
}

type Hello interface {
	SayHello() string
}

func (h *HelloRemoteObject) SayHello() string {
	return h.helloSentence
}

type Calculator interface {
	Sum(a float32, b float32) float32
	Subtract(a float32, b float32) float32
	Multiple(a float32, b float32) float32
	Devide(a float32, b float32) float32
}
