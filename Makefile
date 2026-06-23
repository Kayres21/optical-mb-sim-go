.PHONY: all clean test build run

APP_NAME := simulador
BIN_DIR  := bin

CONFIG ?= files/config.json

all: test build

clean:
	@echo "Cleaning previous builds..."
	@rm -rf $(BIN_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

build: clean
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) .

run: build
	@echo "Running $(APP_NAME) using config=$(CONFIG)..."
	@./$(BIN_DIR)/$(APP_NAME) -config="$(CONFIG)"
