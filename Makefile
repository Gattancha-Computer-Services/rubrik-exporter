.PHONY: build run clean help deps

# Variables
BINARY_NAME=rubrik-exporter
GO=go
GFLAGS=-v
LDFLAGS=-ldflags="-s -w"

help:
	@echo "Available targets:"
	@echo "  make build       - Build the binary"
	@echo "  make run         - Build and run the exporter"
	@echo "  make clean       - Remove the binary"
	@echo "  make deps        - Download dependencies"

deps:
	$(GO) mod download
	$(GO) mod tidy

build: deps
	$(GO) build $(GFLAGS) $(LDFLAGS) -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME) -rubrik.url $(RUBRIK_URL) -rubrik.username $(RUBRIK_USER) -rubrik.password $(RUBRIK_PASSWORD)

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
