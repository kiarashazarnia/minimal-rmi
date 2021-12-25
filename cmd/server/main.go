package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/kiarashazarnia/minimal-rmi/pkg/rmi"
)

var objectsContext = make(map[string]interface{})

func Invoke(any interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	return reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

func handleMethodCall(methodCall rmi.MethodCall) {
	values := Invoke(
		objectsContext[rmi.GenerateKey(methodCall.Target.Name, methodCall.Target.Version)],
		methodCall.MethodName,
		methodCall.Parameters)
	log.Println("method call result:", values)
}

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

	fmt.Println("1")
	rmi.WaitForServer(config.RMI_HOST)
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
