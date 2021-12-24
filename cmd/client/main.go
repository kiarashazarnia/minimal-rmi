package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
)

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

type Hello interface {
	SayHello() string
}

type HelloStub struct {
	remoteAddress string
}

func (h *HelloStub) SayHello() string {
	data = {"dummy": "dummy"}
	body := json.Marshal(data)
	response, err = http.Post(h.remoteAddress, "application/json", body)
	return string(response.Body)
}


func lookup(objectType Type, version int) interface{} {
	return nil
}



func main() {

	hello := lookup(Hello, 1)
	result := hello.SayHello()
	log.remote("Hello object remote call:" + result)
}
