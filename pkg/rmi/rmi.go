package rmi

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Hello interface {
	SayHello() string
	SayHelloTo(name string) string
}

type Calculator interface {
	Sum(a float32, b float32) float32
	Subtract(a float32, b float32) float32
	Multiple(a float32, b float32) float32
	Devide(a float32, b float32) float32
}

type Fibonacci interface {
	Fibonacci(num uint64) uint64
}

type RegisterObjectCommand struct {
	Version       uint
	Name          string
	RemoteAddress string
}

type LookupQuery struct {
	Version uint
	Name    string
}

type LookupResponse struct {
	RemoteAddress string
}

type ServerStub interface {
	Name() string
	Version() uint
}

type ClientStub interface {
	Name() string
	Version() uint
	SetRemoteAddress(remoteAddress string)
	SetVersion(version uint)
}

type MethodCall struct {
	ObjectName    string
	Version       uint
	MethodName    string
	Parameters    string
	HasParameters bool
}

func WaitForServer(host string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		out, _ := exec.CommandContext(ctx, "echo", "hello").Output()
		if string(out) == "hello" || ctx.Err() != nil {
			break
		}

		timeout := time.Second
		conn, err := net.DialTimeout("tcp", host, timeout)
		if err != nil {
			log.Println("Waiting for server", err)
		}
		if conn != nil {
			defer conn.Close()
			fmt.Println("Server is responsive:", host)
			break
		}

	}
}

func GenerateKey(name string, version uint) string {
	key := fmt.Sprintf("%s:%d", name, version)
	// log.Println("generated key:", key)
	return key
}

func RMIUrl(address string) string {
	return fmt.Sprintf("http://%s/remote", address)
}

func DecodeArguments(parameters string) []interface{} {

	// log.Println("decoding args:", parameters)

	params := strings.Split(parameters, "|")
	arguments := make([]interface{}, len(params))

	// log.Println("parsed parameters:", params)

	for i, parameter := range params {
		arguments[i] = parameter
		uintValue, err := strconv.ParseUint(parameter, 10, 64)
		if err == nil {
			arguments[i] = uintValue
			continue
		}
		intValue, err := strconv.Atoi(parameter)
		if err == nil {
			arguments[i] = intValue
			continue
		}
		floatValue, err := strconv.ParseFloat(parameter, 32)
		if err == nil {
			arguments[i] = floatValue
			continue
		}
		boolValue, err := strconv.ParseBool(parameter)
		if err == nil {
			arguments[i] = boolValue
			continue
		}
	}
	// log.Println("decoded args:", arguments)
	return arguments
}

func EncodeArguments(args ...interface{}) string {
	buffer := bytes.NewBufferString("")
	for _, arg := range args {
		s := fmt.Sprintf("%v", arg)
		buffer.WriteString(s)
	}
	res := buffer.String()
	// log.Println("encoded args:", res)
	return res
}
