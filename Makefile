BINARY := bin/unix-socket-path-limit

.DEFAULT_GOAL := help
.PHONY: help demo build clean

help:
	@echo "Usage:"
	@echo "  make demo    - Run the demo program"
	@echo "  make build   - Build the binary executable"
	@echo "  make help    - Show this help message"

demo:
	@echo "Running demo..."
	@go run .
	@echo "Demo completed."

build:
	@echo "Building binary..."
	@go build -o $(BINARY) .
	@echo "Binary built at $(BINARY)"

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY)
	@echo "Clean completed."