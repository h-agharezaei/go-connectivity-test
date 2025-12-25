.PHONY: build run clean test deps

# Build the application
build:
	go build -o bin/connectivity-test ./cmd

# Run the application
run: build
	./bin/connectivity-test

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy
