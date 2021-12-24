package main

import (
	"fmt"
	"net/http"
)

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	// implement objects methods
	// implement parsing and performing incoming command
	// initialize objects
	// register them to registry

	var hello Hello = new(HelloRemoteObject)

	register(hello)

	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8090", nil)
}

func register(hello interface{}) {
	data := {
		"version": 	1,
		"name": "hello"
	}
	jsonData, err = json.Marshal(data)
	resp, err = http.Post("http://localhost:8080", "application/json", bytes.newBuffer(jsonData))
}

type HelloRemoteObject struct {
	helloSentence string
}

type Hello interface {
	SayHello() string
}

func (h *HelloRemoteObject) SayHello() string {
	return h.helloSentence
}

type Calculator interface {
	Sum(a float32, b float32) float32
	Subtract(a float32, b float32) float32
	Multiple(a float32, b float32) float32
	Devide(a float32, b float32) float32
}
