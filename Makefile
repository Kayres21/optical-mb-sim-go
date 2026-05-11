.PHONY: all clean test build run

APP_NAME := simulador
BIN_DIR  := bin

# ── Simulation parameters (override with: make run BANDS=4 LAMBDA=100) ────────
NETWORK    ?= files/networks/UKNet_BDM.json
ROUTES     ?= files/routes/UKNet_routes.json
CAPACITIES ?= files/capacities/capacities.json
BITRATE    ?= files/bitrate/bitrate.json
LAMBDA     ?= 50
MU         ?= 1
BANDS      ?= 1
GOAL       ?= 1e8
LOGS       ?= true
LEGACY     ?= false

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
	@echo "Running $(APP_NAME) (λ=$(LAMBDA) μ=$(MU) bands=$(BANDS) goal=$(GOAL) legacy=$(LEGACY))..."
	@./$(BIN_DIR)/$(APP_NAME) \
		-network="$(NETWORK)" \
		-routes="$(ROUTES)" \
		-capacities="$(CAPACITIES)" \
		-bitrate="$(BITRATE)" \
		-lambda=$(LAMBDA) \
		-mu=$(MU) \
		-bands=$(BANDS) \
		-goal=$(GOAL) \
		-logs=$(LOGS) \
		-legacy=$(LEGACY)
