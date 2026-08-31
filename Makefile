.PHONY: all clean test build run

APP_NAME := simulador
BIN_DIR  := bin

CONFIG ?= files/config.json
LOGS ?= true
EVENTS_CSV ?=
SWEEP ?= false
LAMBDA_START ?= 500
LAMBDA_END ?= 1500
LAMBDA_STEP ?= 50

RUN_FLAGS = -config="$(CONFIG)" \
	-logs="$(LOGS)" \
	-sweep="$(SWEEP)" \
	-lambda-start="$(LAMBDA_START)" \
	-lambda-end="$(LAMBDA_END)" \
	-lambda-step="$(LAMBDA_STEP)" \
	$(if $(EVENTS_CSV),-events-csv="$(EVENTS_CSV)",)

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
	@./$(BIN_DIR)/$(APP_NAME) $(RUN_FLAGS)
