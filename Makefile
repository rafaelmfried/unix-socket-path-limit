BINARY := bin/unix-socket-path-limit

.DEFAULT_GOAL := help
.PHONY: help demo

help:
	@echo "Usage:"
	@echo "  make demo    - Run the demo program"
	@echo "  make help    - Show this help message"

demo:
	@echo "Running demo..."
	@go run .
	@echo "Demo completed."