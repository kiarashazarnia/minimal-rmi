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

type HelloStub struct {
	remoteAddress string
}

func (h *HelloStub) SayHello() string {
	log.Println("client stub saying hello")

	methodCall := rmi.MethodCall{
		ObjectName: "Hello",
		Version:    1,
		MethodName: "SayHello",
		Parameters: "",
	}
	body, _ := json.Marshal(methodCall)
	requestBody := bytes.NewBuffer(body)
	url := rmi.RMIUrl(h.remoteAddress)
	log.Println("sending request:", requestBody, " address:", url)
	response, err := http.Post(url, "application/json", requestBody)
	log.Println("response:", response, err)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	log.Println("response:", string(responseBody))
	return string(responseBody)
}

func (h *HelloStub) Name() string {
	return "Hello"
}

func (h *HelloStub) Version() uint {
	return 1
}

func (h *HelloStub) SetRemoteAddress(address string) {
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
	stubRegistry[rmi.GenerateKey("Hello", 1)] = reflect.TypeOf(HelloStub{})
}

func makeInstance(name string, remoteAdreess string) interface{} {
	log.Println("making stub type:", name, " object:", stubRegistry[name])
	v := reflect.New(stubRegistry[name]).Elem()
	log.Println("object value:", &v)
	var helloStub HelloStub = v.Interface().(HelloStub)
	log.Println("value interface:", reflect.TypeOf(helloStub))
	var stub rmi.StubObject = &helloStub
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
	result := hello.SayHello()
	log.Print("rmi.Hello object remote call:" + result)
	// ExecuteCommand()
}
