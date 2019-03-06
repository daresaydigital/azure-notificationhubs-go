default: test

test:
	go test -vet -v -cover -race ./...

build:
	go build ./...
