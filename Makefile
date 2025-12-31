run:
	go run cmd/main.go

build:
	go build -o bin/$(BINARY_NAME) cmd/main.go