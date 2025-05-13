APP_NAME=ShutterSync
BINARY_NAME=shuttersync
VERSION=1.0.0
BIN_DIR=bin
GO_BUILD=go build

.PHONY: build run clean

build:
	@echo "Building binary..."
	@mkdir -p $(BIN_DIR)
	$(GO_BUILD) -o $(BIN_DIR)/$(BINARY_NAME) cmd/*.go

run: build
	@echo "Running $(APP_NAME)..."
	./$(BIN_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
