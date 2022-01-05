package rmi

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

type Hello interface {
	SayHello() string
}

type Calculator interface {
	Sum(a float32, b float32) float32
	Subtract(a float32, b float32) float32
	Multiple(a float32, b float32) float32
	Devide(a float32, b float32) float32
}

type Salam struct {
	Name string
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

type StubObject interface {
	Name() string
	Version() uint
	SetRemoteAddress(remoteAddress string)
}

type MethodCall struct {
	ObjectName string
	Version    uint
	MethodName string
	Parameters string
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
	log.Println("generated key:", key)
	return key
}

func RMIUrl(address string) string {
	return fmt.Sprintf("http://%s/remote", address)
}
