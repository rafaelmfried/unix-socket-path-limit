BINARY := bin/unix-socket-path-limit

.DEFAULT_GOAL := help
.PHONY: help demo build clean

help:
	@echo "Usage:"
	@echo "  make demo    - Run the demo program"
	@echo "  make test    - Run tests"
	@echo "  make build   - Build the binary executable"
	@echo "  make clean   - Remove the built binary"
	@echo "  make help    - Show this help message"

demo:
	@echo "Running demo..."
	@go run .
	@echo "Demo completed."

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests completed."

build:
	@echo "Building binary..."
	@go build -o $(BINARY) .
	@echo "Binary built at $(BINARY)"

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY)
	@echo "Clean completed."