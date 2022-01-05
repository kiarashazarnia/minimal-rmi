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
	./parallel_commands.sh ./server.out ./client.out ./rmi.out

.PHONY: clean
clean:
	rm -rf *.out
	go mod vendor

.PHONY: test
test:
	go test -mod=vendor -v ./... -coverprofile cover.out
