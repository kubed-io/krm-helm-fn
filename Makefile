.PHONY: build test docker-build docker-push

# Default target
all: build

# Build the function
build:
	go build -o bin/function ./cmd/function

# Run tests
test:
	go test -v ./...

# Build Docker image
docker-build:
	docker build -t kubed/krm-helm-fn:latest .

# Push Docker image
docker-push:
	docker push kubed/krm-helm-fn:latest

# Clean
clean:
	rm -rf bin/

# Run with example
run-example:
	kubectl kustomize examples/inflate

# Install dependencies
deps:
	go mod download