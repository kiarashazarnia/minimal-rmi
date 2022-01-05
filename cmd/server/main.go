package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/kiarashazarnia/minimal-rmi/pkg/rmi"
)

var objectsContext = make(map[string]interface{})

func Invoke(any interface{}, name string, args ...interface{}) []reflect.Value {
	log.Println("interface:", any, "name:", name, "args:", args, "args len:", len(args))
	inputs := make([]reflect.Value, len(args))
	log.Println("inputs:", inputs, "len:", len(inputs))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	return reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

func handleMethodCall(methodCall rmi.MethodCall) {
	log.Println("invoking method:", methodCall)
	var values []reflect.Value
	if methodCall.HasParameters {
		values = Invoke(
			objectsContext[rmi.GenerateKey(methodCall.ObjectName, methodCall.Version)],
			methodCall.MethodName,
			methodCall.Parameters)
	} else {
		values = Invoke(
			objectsContext[rmi.GenerateKey(methodCall.ObjectName, methodCall.Version)],
			methodCall.MethodName)
	}
	log.Println("method call result:", values)
}

func remote(w http.ResponseWriter, req *http.Request) {

	log.Println("handing remote method invocation")

	decoder := json.NewDecoder(req.Body)
	var methodCall rmi.MethodCall
	err := decoder.Decode(&methodCall)
	if err != nil {
		log.Println("decoding error:", err)
		return
	}
	log.Println("decoded RMI request:", methodCall)
	handleMethodCall(methodCall)
}

var config = rmi.LoadConfig()

func main() {

	rmi.WaitForServer(config.RMI_HOST)
	http.HandleFunc("/remote", remote)

	// we instantiate HelloRemoteObject but save it in an Hello type variable
	var hello rmi.Hello = HelloRemoteObject{
		helloSentence: "Hello World",
	}
	register(hello.(rmi.ServerStub))
	log.Println("running remote server on:", config.REMOTE_HOST)
	err := http.ListenAndServe(config.REMOTE_HOST, nil)
	log.Println("error occured:", err)
}

func register(object rmi.ServerStub) bool {
	objectsContext[rmi.GenerateKey(object.Name(), object.Version())] = object
	data := rmi.RegisterObjectCommand{
		Version:       object.Version(),
		Name:          object.Name(),
		RemoteAddress: config.REMOTE_HOST,
	}
	jsonData, _ := json.Marshal(data)
	log.Println("registering remote object:", data)
	url := "http://" + config.RMI_HOST + "/register"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	log.Println("registeration result and error:", response, err)
	return response.StatusCode == http.StatusOK
}

type HelloRemoteObject struct {
	helloSentence string
}

func (h HelloRemoteObject) Name() string {
	return "Hello"
}

func (h HelloRemoteObject) Version() uint {
	return 1
}

func (h HelloRemoteObject) SayHello() string {
	return h.helloSentence
}

type Calculator interface {
	Sum(a float32, b float32) float32
	Subtract(a float32, b float32) float32
	Multiple(a float32, b float32) float32
	Devide(a float32, b float32) float32
}
