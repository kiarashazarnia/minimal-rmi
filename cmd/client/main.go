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
	stubRegistry[rmi.GenerateKey("Factorial", 1)] = reflect.TypeOf(FactorialClientStub{})
	stubRegistry[rmi.GenerateKey("Factorial", 2)] = reflect.TypeOf(FactorialClientStub{})
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
	case FactorialClientStub:
		var factorialStub FactorialClientStub = v.Interface().(FactorialClientStub)
		log.Println("value interface:", reflect.TypeOf(factorialStub))
		stub = &factorialStub
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

	factorialStubIface1 := lookup("Factorial", 1)
	var factorialStub1 rmi.Factorial = factorialStubIface1.(rmi.Factorial)
	factResult1 := factorialStub1.Factorial(20)
	log.Printf("factroial %d!=%d", 20, factResult1)

	factorialStubIface2 := lookup("Factorial", 2)
	var factorialStub2 rmi.Factorial = factorialStubIface2.(rmi.Factorial)
	factResult2 := factorialStub2.Factorial(20)
	log.Printf("factroial %d!=%d", 20, factResult2)
}
