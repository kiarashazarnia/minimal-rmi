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

func lookup(objectType reflect.Type, version int) interface{} {
	salam := rmi.Salam{Name: "salam"}
	fmt.Println(salam)
	return HelloStub{
		name:          "name",
		version:       version,
		remoteAddress: "localhost",
	}
}

func main() {
	var hello rmi.Hello
	rmi.WaitForServer(config.RMI_HOST)
	hello = lookup(reflect.TypeOf(hello), 1).(rmi.Hello)
	result := hello.SayHello()
	log.Print("rmi.Hello object remote call:" + result)
}
