default: all

all: test lint vet build

test:
	go test -vet -v -cover -race ./...

build:
	go build ./...

coverage:
	go test -v -cover -coverprofile coverage.out -race ./...

lint:
	golint -set_exit_status ./...

vet:
	go vet ./...
