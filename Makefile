.PHONY: build
build:
	go build -a -mod=vendor -o client.out ./cmd/client/main.go
	go build -a -mod=vendor -o server.out ./cmd/server/main.go
	go build -a -mod=vendor -o rmi.out ./cmd/rmi/main.go

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
