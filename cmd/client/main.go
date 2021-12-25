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
	version       int
	remoteAddress string
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (h *HelloStub) SayHello() string {
	body, _ := json.Marshal(h)
	requestBody := bytes.NewBuffer(body)
	response, _ := http.Post(h.remoteAddress, "application/json", requestBody)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	return string(responseBody)
}

func lookup(name string, version uint) interface{} {

	data := rmi.LookupCommand{
		Version: version,
		Name:    name,
	}
	jsonData, _ := json.Marshal(data)
	log.Println("looking up remote object:", data)
	url := "http://" + config.RMI_HOST + "/lookup"
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	log.Println("looking up result and error:", response, err)
	// todo deserialize hello stub
	return nil
}

func main() {
	rmi.WaitForServer(config.RMI_HOST)
	var hello rmi.Hello = lookup("<main.HelloRemoteObject Value>", 1).(rmi.Hello)
	result := hello.SayHello()
	log.Print("rmi.Hello object remote call:" + result)
}
