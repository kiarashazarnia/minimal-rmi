.PHONY: build
build: fmt
	go build -mod=modules -a -o client ./cmd/main.go
	go build -mod=modules -a -o server ./cmd/main.go
	go build -mod=modules -a -o rmi ./cmd/main.go

.PHONY: run
run:
	./rmi
	./server
	./client

.PHONY: test
test:
	go test -mod=modules -v ./... -coverprofile cover.out

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: fmt-check
fmt-check: fmt
	git diff-index --quiet HEAD