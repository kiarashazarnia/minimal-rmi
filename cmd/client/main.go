package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/kiarashazarnia/minimal-rmi/pkg/rmi"
)

var config = rmi.LoadConfig()

var stubRegistry = make(map[string]reflect.Type)

func lookupStub(name string, version uint, lookupResult rmi.LookupResponse) interface{} {
	return makeInstance(name, version, lookupResult.RemoteAddress)
}

func initStubRegistry() {
	stubRegistry[rmi.GenerateKey("Hello", 1)] = reflect.TypeOf(HelloClientStub{})
	stubRegistry[rmi.GenerateKey("Hello", 2)] = reflect.TypeOf(HelloClientStub{})
	stubRegistry[rmi.GenerateKey("Fibonacci", 1)] = reflect.TypeOf(FibonacciClientStub{})
	stubRegistry[rmi.GenerateKey("Fibonacci", 2)] = reflect.TypeOf(FibonacciClientStub{})
}

func makeInstance(name string, version uint, remoteAdreess string) interface{} {
	stubKey := rmi.GenerateKey(name, version)
	log.Println("making stub type:", stubKey, " object:", stubRegistry[stubKey])
	v := reflect.New(stubRegistry[stubKey]).Elem()
	log.Println("object value:", &v)
	var stub rmi.ClientStub = nil

	switch stubType := v.Interface().(type) {
	case HelloClientStub:
		var helloStub HelloClientStub = v.Interface().(HelloClientStub)
		log.Println("value interface:", reflect.TypeOf(helloStub))
		stub = &helloStub
	case FibonacciClientStub:
		var fibonacciStub FibonacciClientStub = v.Interface().(FibonacciClientStub)
		log.Println("value interface:", reflect.TypeOf(fibonacciStub))
		stub = &fibonacciStub
	default:
		log.Println("error type not registerd:", stubType)
		return nil
	}
	log.Println("conversion result:", stub)
	stub.SetRemoteAddress(remoteAdreess)
	stub.SetVersion(version)
	return stub
}

func lookup(name string, version uint) interface{} {
	data := rmi.LookupQuery{
		Version: version,
		Name:    name,
	}
	jsonData, _ := json.Marshal(data)
	log.Println("looking up remote object:", data)
	url := "http://" + config.RMI_HOST + "/lookup"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	log.Println("looking up result:", response, err)
	var lookupResult rmi.LookupResponse = rmi.LookupResponse{}
	json.NewDecoder(response.Body).Decode(&lookupResult)
	return lookupStub(name, version, lookupResult)
}

func main() {
	initStubRegistry()
	rmi.WaitForServer(config.RMI_HOST)

	helloStubIface1 := lookup("Hello", 1)
	var hello1 rmi.Hello = helloStubIface1.(rmi.Hello)
	hello1.SayHello()
	hello1.SayHelloTo("Amir")

	helloStubIface2 := lookup("Hello", 2)
	var hello2 rmi.Hello = helloStubIface2.(rmi.Hello)
	hello2.SayHello()
	hello2.SayHelloTo("Amir")

	fibonacciStubIface1 := lookup("Fibonacci", 1)
	var fibonacciStub1 rmi.Fibonacci = fibonacciStubIface1.(rmi.Fibonacci)
	fiboResult1 := fibonacciStub1.Fibonacci(50)
	log.Printf("fiboroial %d!=%d", 50, fiboResult1)

	fibonacciStubIface2 := lookup("Fibonacci", 2)
	var fibonacciStub2 rmi.Fibonacci = fibonacciStubIface2.(rmi.Fibonacci)
	fiboResult2 := fibonacciStub2.Fibonacci(50)
	log.Printf("fiboroial %d!=%d", 50, fiboResult2)
}
