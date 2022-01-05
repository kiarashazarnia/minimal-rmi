package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kiarashazarnia/minimal-rmi/pkg/rmi"
)

type FibonacciClientStub struct {
	version       uint
	remoteAddress string
}

func (fibonacciStub *FibonacciClientStub) Fibonacci(num uint64) uint64 {
	start := time.Now()
	methodCall := rmi.MethodCall{
		ObjectName:    fibonacciStub.Name(),
		Version:       fibonacciStub.Version(),
		MethodName:    "Fibonacci",
		Parameters:    rmi.EncodeArguments(num),
		HasParameters: true,
	}
	log.Println("client stub calling method:", methodCall)
	body, _ := json.Marshal(methodCall)
	requestBody := bytes.NewBuffer(body)
	url := rmi.RMIUrl(fibonacciStub.remoteAddress)
	log.Println("sending request:", requestBody, " address:", url)
	response, err := http.Post(url, "application/json", requestBody)
	log.Println("response:", response, err)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	log.Println("response:", string(responseBody))
	result := rmi.DecodeArguments(string(responseBody))[0].(uint64)
	log.Println("RMI result:", result, err)
	elapsed := time.Since(start)
	log.Printf("**** Fibonacci version:%d took %s ******\n", methodCall.Version, elapsed)
	return result
}

func (f *FibonacciClientStub) Name() string {
	return "Fibonacci"
}

func (f *FibonacciClientStub) Version() uint {
	return f.version
}

func (f *FibonacciClientStub) SetRemoteAddress(address string) {
	f.remoteAddress = address
}

func (f *FibonacciClientStub) SetVersion(version uint) {
	f.version = version
}

type HelloClientStub struct {
	remoteAddress string
	version       uint
}

func (h *HelloClientStub) SayHello() string {
	log.Println("client stub saying hello")

	methodCall := rmi.MethodCall{
		ObjectName:    h.Name(),
		Version:       h.Version(),
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
		ObjectName:    h.Name(),
		Version:       h.Version(),
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
	return h.version
}

func (h *HelloClientStub) SetRemoteAddress(address string) {
	h.remoteAddress = address
}

func (h *HelloClientStub) SetVersion(version uint) {
	h.version = version
}
