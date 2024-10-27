# Makefile

# Variables
BINARY_NAME := shuttersync

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GORUN := $(GOCMD) run
GOGET := $(GOCMD) get

# Directories
CMD_DIR := cmd
BIN_DIR := bin

.PHONY: all run build test clean deps install

all: build

build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/*.go

run: build
	$(GORUN) $(CMD_DIR)/*.go

test:
	$(GOTEST) ./...

clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)/$(BINARY_NAME)

install: build
	mv $(BIN_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
