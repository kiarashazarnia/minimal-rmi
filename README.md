# Mimial RMI

###### A minimal implementation to understand Remote Method Invocation concepts better.

This codebase is submitted as a computer assignment for distributed systems course offered by Prof. A. Kalbasi at Tehran Polytechnic.


## Run


## Description 

```golang


//************** rmi server *****************

type Hello interface {
	SayHello() string
}

// ***************** remote object server *****************

type HelloRemoteObject struct {
	helloSentence string
}

func (h HelloRemoteObject) SayHello() string {
	return h.helloSentence
}


// ***************** client *****************

type HelloStub struct {
	name          string
	version       int
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

func (h *HelloStub) SetRemoteAddress(remoteAddress string) {
	h.remoteAddress = remoteAddress
}


var hello rmi.Hello = lookup("<rmi.Hello Value>", 1).(rmi.Hello)
result := hello.SayHello()

```

## References
These material are used to implement this code:

1. https://talks.golang.org/2013/distsys.slide#47
2. https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
3. https://dev.to/jinxankit/go-project-structure-and-guidelines-4ccm
4. https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go
5. https://stackoverflow.com/questions/8103617/call-a-struct-and-its-method-by-name-in-go
6. https://stackoverflow.com/questions/10909685/run-parallel-multiple-commands-at-once-in-the-same-terminal 
7. 
8. 
