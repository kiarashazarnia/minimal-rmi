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
	name          string
	version       uint
	remoteAddress string
}

func (h *HelloStub) SayHello() string {
	body, _ := json.Marshal(h)
	requestBody := bytes.NewBuffer(body)
	response, _ := http.Post(h.remoteAddress, "application/json", requestBody)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	return string(responseBody)
}

func (h *HelloStub) Name() string {
	return h.name
}

func (h *HelloStub) Version() uint {
	return h.version
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
	myTypes := []interface{}{HelloStub{}}
	for _, v := range myTypes {
		var stub rmi.StubObject = v.(rmi.StubObject)
		classKey := rmi.GenerateKey(stub.Name(), stub.Version())
		stubRegistry[classKey] = reflect.TypeOf(stub)
	}
}

func makeInstance(name string, remoteAdreess string) interface{} {
	v := reflect.New(stubRegistry[name]).Elem()
	var vInterface interface{} = v
	vInterface.(rmi.StubObject).SetRemoteAddress(remoteAdreess)
	return v.Interface()
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
	log.Println("looking up result and error:", response, err)
	var lookupResult rmi.LookupResponse = rmi.LookupResponse{}
	json.NewDecoder(response.Body).Decode(&lookupResult)
	return lookupStub(name, version, lookupResult)
}

func main() {
	initStubRegistry()
	rmi.WaitForServer(config.RMI_HOST)
	var hello rmi.Hello = lookup("<rmi.Hello Value>", 1).(rmi.Hello)
	result := hello.SayHello()
	log.Print("rmi.Hello object remote call:" + result)
	// ExecuteCommand()
}
