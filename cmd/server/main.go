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
	log.Println("interface:", any, "name:", name, "args:", args, "args len:", len(args))
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	log.Println("inputs:", inputs, "len:", len(inputs))
	return reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

func handleMethodCall(methodCall rmi.MethodCall) string {
	log.Println("invoking method:", methodCall)
	var values []reflect.Value
	if methodCall.HasParameters {
		values = Invoke(
			objectsContext[rmi.GenerateKey(methodCall.ObjectName, methodCall.Version)],
			methodCall.MethodName,
			rmi.DecodeArguments(methodCall.Parameters)...)
	} else {
		values = Invoke(
			objectsContext[rmi.GenerateKey(methodCall.ObjectName, methodCall.Version)],
			methodCall.MethodName)
	}
	log.Println("method call result:", values)

	result := rmi.EncodeArguments(convertToIfaces(values)...)
	return result
}

func convertToIfaces(vals []reflect.Value) []interface{} {
	args := make([]interface{}, len(vals))
	for i, v := range vals {
		args[i] = v.Interface()
	}
	return args
}

func remote(w http.ResponseWriter, req *http.Request) {

	log.Println("handing remote method invocation")

	decoder := json.NewDecoder(req.Body)
	var methodCall rmi.MethodCall
	err := decoder.Decode(&methodCall)
	if err != nil {
		log.Println("decoding error:", err)
		w.WriteHeader(500)
		return
	}
	log.Println("decoded RMI request:", methodCall)
	result := handleMethodCall(methodCall)
	w.WriteHeader(200)
	w.Write([]byte(result))
	return
}

var config = rmi.LoadConfig()

func main() {

	rmi.WaitForServer(config.RMI_HOST)
	http.HandleFunc("/remote", remote)
	initServerStubs()
	log.Println("running remote server on:", config.REMOTE_HOST)
	err := http.ListenAndServe(config.REMOTE_HOST, nil)
	log.Println("error occured:", err)
}

func initServerStubs() {

	// we instantiate HelloRemoteObject but save it in an Hello type variable
	var hello1 rmi.Hello = HelloServerStub{
		helloSentence: "Hello RMI World version 1!",
		version:       1,
	}
	var hello2 rmi.Hello = HelloServerStub{
		helloSentence: "Hello RMI World version 2!",
		version:       2,
	}
	register(hello1.(rmi.ServerStub))
	register(hello2.(rmi.ServerStub))

	var factorial1 rmi.Factorial = RecursiveFactorialServerStub{}
	var factorial2 rmi.Factorial = DynamicFactorialServerStub{}
	register(factorial1.(rmi.ServerStub))
	register(factorial2.(rmi.ServerStub))
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

type HelloServerStub struct {
	helloSentence string
	version       uint
}

func (h HelloServerStub) Name() string {
	return "Hello"
}

func (h HelloServerStub) Version() uint {
	return h.version
}

func (h HelloServerStub) SayHello() string {
	return h.helloSentence
}

func (h HelloServerStub) SayHelloTo(name string) string {
	return fmt.Sprintf("Hello dear %s!", name)
}

type RecursiveFactorialServerStub struct {
}

func (s RecursiveFactorialServerStub) Version() uint {
	return 1
}

func (s RecursiveFactorialServerStub) Name() string {
	return "Factorial"
}

func (s RecursiveFactorialServerStub) Factorial(num uint64) uint64 {
	if num <= 1 {
		return 1
	}
	return num * s.Factorial(num-1)
}

type DynamicFactorialServerStub struct {
}

func (s DynamicFactorialServerStub) Version() uint {
	return 2
}

func (s DynamicFactorialServerStub) Name() string {
	return "Factorial"
}

func (s DynamicFactorialServerStub) Factorial(num uint64) uint64 {
	factVal := uint64(1)
	for i := uint64(1); i <= num; i++ {
		factVal *= i
	}
	return factVal
}
