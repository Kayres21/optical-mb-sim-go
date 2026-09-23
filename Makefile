.PHONY: all clean test build run

APP_NAME := simulador
BIN_DIR  := bin

CONFIG ?= files/config.json
LOGS ?= true
EVENTS_CSV ?=
SWEEP ?=
LAMBDA_START ?=
LAMBDA_END ?=
LAMBDA_STEP ?=

RUN_FLAGS = -config="$(CONFIG)" \
	-logs="$(LOGS)" \
	$(if $(SWEEP),-sweep="$(SWEEP)",) \
	$(if $(LAMBDA_START),-lambda-start="$(LAMBDA_START)",) \
	$(if $(LAMBDA_END),-lambda-end="$(LAMBDA_END)",) \
	$(if $(LAMBDA_STEP),-lambda-step="$(LAMBDA_STEP)",) \
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
