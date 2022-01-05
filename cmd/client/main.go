package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/kiarashazarnia/minimal-rmi/pkg/rmi"
)

var config = rmi.LoadConfig()

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

type HelloClientStub struct {
	remoteAddress string
}

func (h *HelloClientStub) SayHello() string {
	log.Println("client stub saying hello")

	methodCall := rmi.MethodCall{
		ObjectName:    "Hello",
		Version:       1,
		MethodName:    "SayHello",
		Parameters:    "",
		HasParameters: false,
	}
	body, _ := json.Marshal(methodCall)
	requestBody := bytes.NewBuffer(body)
	url := rmi.RMIUrl(h.remoteAddress)
	log.Println("sending request:", requestBody, " address:", url)
	response, err := http.Post(url, "application/json", requestBody)
	log.Println("response:", response, err)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	result := string(responseBody)
	log.Println("RMI result:", result)
	return result
}

func (h *HelloClientStub) SayHelloTo(name string) string {
	methodCall := rmi.MethodCall{
		ObjectName:    "Hello",
		Version:       1,
		MethodName:    "SayHelloTo",
		Parameters:    rmi.EncodeArguments(name),
		HasParameters: true,
	}
	log.Println("client stub calling method:", methodCall)
	body, _ := json.Marshal(methodCall)
	requestBody := bytes.NewBuffer(body)
	url := rmi.RMIUrl(h.remoteAddress)
	log.Println("sending request:", requestBody, " address:", url)
	response, err := http.Post(url, "application/json", requestBody)
	log.Println("response:", response, err)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	log.Println("response:", string(responseBody))
	result := string(responseBody)
	log.Println("RMI result:", result)
	return result
}

func (h *HelloClientStub) Name() string {
	return "Hello"
}

func (h *HelloClientStub) Version() uint {
	return 1
}

func (h *HelloClientStub) SetRemoteAddress(address string) {
	h.remoteAddress = address
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

var stubRegistry = make(map[string]reflect.Type)

func lookupStub(name string, version uint, lookupResult rmi.LookupResponse) interface{} {
	stubKey := rmi.GenerateKey(name, version)
	return makeInstance(stubKey, lookupResult.RemoteAddress)
}

func initStubRegistry() {
	stubRegistry[rmi.GenerateKey("Hello", 1)] = reflect.TypeOf(HelloClientStub{})
	stubRegistry[rmi.GenerateKey("Hello", 2)] = reflect.TypeOf(HelloClientStub{})
	stubRegistry[rmi.GenerateKey("Factorial", 1)] = reflect.TypeOf(FactorialClientStub{})
	stubRegistry[rmi.GenerateKey("Factorial", 2)] = reflect.TypeOf(FactorialClientStub{})
}

func makeInstance(name string, remoteAdreess string) interface{} {
	log.Println("making stub type:", name, " object:", stubRegistry[name])
	v := reflect.New(stubRegistry[name]).Elem()
	log.Println("object value:", &v)
	var stub rmi.StubObject = nil

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
	lookedUp := lookup("Hello", 1)
	log.Print("looked up:", reflect.TypeOf(lookedUp), lookedUp)
	var hello rmi.Hello = lookedUp.(rmi.Hello)
	hello.SayHello()
	hello.SayHelloTo("Amir")

	factorialStubInterface := lookup("Factorial", 1)
	log.Print("looked up:", reflect.TypeOf(factorialStubInterface), factorialStubInterface)
	var factorialStub1 rmi.Factorial = factorialStubInterface.(rmi.Factorial)
	factResult := factorialStub1.Factorial(10)
	log.Printf("factroial %d=%d", 10, factResult)
}

type FactorialClientStub struct {
	remoteAddress string
}

func (h *FactorialClientStub) Factorial(num uint64) uint64 {
	methodCall := rmi.MethodCall{
		ObjectName:    h.Name(),
		Version:       h.Version(),
		MethodName:    "Factorial",
		Parameters:    rmi.EncodeArguments(num),
		HasParameters: true,
	}
	log.Println("client stub calling method:", methodCall)
	body, _ := json.Marshal(methodCall)
	requestBody := bytes.NewBuffer(body)
	url := rmi.RMIUrl(h.remoteAddress)
	log.Println("sending request:", requestBody, " address:", url)
	response, err := http.Post(url, "application/json", requestBody)
	log.Println("response:", response, err)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	log.Println("response:", string(responseBody))
	result := rmi.DecodeArguments(string(responseBody))[0].(uint64)
	log.Println("RMI result:", result, err)
	return result
}

func (h *FactorialClientStub) Name() string {
	return "Factorial"
}

func (h *FactorialClientStub) Version() uint {
	return 1
}

func (h *FactorialClientStub) SetRemoteAddress(address string) {
	h.remoteAddress = address
}
