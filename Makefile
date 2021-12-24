.PHONY: build
build: client server rmi


.PHONY: client
client: ./cmd/client/main.go
	go build -a -mod=vendor -o client.out ./cmd/client/*.go

.PHONY: server
server: ./cmd/server/main.go
	go build -a -mod=vendor -o server.out ./cmd/server/*.go

.PHONY: rmi
rmi: ./cmd/rmi/main.go
	go build -a -mod=vendor -o rmi.out ./cmd/rmi/*.go

.PHONY: run
run:
	cmd/rmi
	cmd/server
	cmd/client

.PHONY: clean
clean:
	rm -rf *.out

.PHONY: test
test:
	go test -mod=vendor -v ./... -coverprofile cover.out
