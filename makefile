# Variables
BINARY_NAME = udp_server
SOURCE_DIR = server
GO_FILES = $(SOURCE_DIR)/*.go

# Targets
all: build

# Build the Go application
build:
	@echo "Building the binary..."
	go build -o $(BINARY_NAME) $(GO_FILES)

# Run the server
run: build
	@echo "Running the server..."
	./$(BINARY_NAME)

# Clean up the binary
clean:
	@echo "Cleaning up..."
	rm -force .\$(BINARY_NAME)

# Install Go dependencies (if any)
deps:
	@echo "Installing dependencies..."
	go mod tidy