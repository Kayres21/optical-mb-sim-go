# Optical Multiband Simulator

This simulator is a performance evaluation tool for optical multiband networks, based on the "Flex Net Sim" project. It evaluates network blocking probability under different arrival rates ($\lambda$), service rates ($\mu$), and frequency band configurations.

## Prerequisites

- **Go**: Version 1.24 or higher is required.
- **Make**: (Optional) For using the provided automation commands.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Kayres21/optical-mb-sim-go.git
   cd optical-mb-sim-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Simulator

The recommended way to run the simulator is using the Makefile.

### Using Make (Recommended)

Run the simulation using the default config file:
```bash
make run
```

Run with a custom config file:
```bash
make run CONFIG=files/config-legacy.json
```

Override simulator flags through Make variables:
```bash
make run CONFIG=files/config.json LOGS=false EVENTS_CSV=files/sim_events.csv
make run SWEEP=true LAMBDA_START=500 LAMBDA_END=1500 LAMBDA_STEP=50
```

### Using Go directly

You can run the binary directly with the supported flags:
```bash
go run main.go -config files/config.json
```

Optional flags:
```bash
go run main.go -config files/config.json -logs=false
go run main.go -config files/config.json -events-csv files/sim_events.csv
go run main.go -config files/config.json -sweep=true -lambda-start 500 -lambda-end 1500 -lambda-step 50
```

Available CLI flags:

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `files/config.json` | Path to the JSON configuration file |
| `-logs` | `true` | Enable progress logging |
| `-events-csv` | empty | Path to save the generated events CSV; overrides `events_csv` when provided |
| `-sweep` | `false` | Run multiple simulations across a lambda range |
| `-lambda-start` | `500` | Start lambda for a sweep |
| `-lambda-end` | `1500` | End lambda for a sweep |
| `-lambda-step` | `50` | Lambda increment for a sweep; must be greater than zero |

## Configuration Parameters

Configuration is provided via JSON. The default file is `files/config.json`:

```json
{
    "network": "files/networks/UKNet_BDM.json",
    "routes": "files/routes/UKNet_routes.json",
    "capacities": "files/capacities/capacities.json",
    "bitrate": "files/bitrate/bitrate.json",
    "lambda": 600,
    "mu": 1,
    "bands": 4,
    "goal": 1e7,
    "logs": true,
    "legacy": false,
    "defrag_mode": "none",
    "sweep": true,
    "lambda_start": 500,
    "lambda_end": 1000,
    "lambda_step": 100
}
```

| Key | Type | Default in files/config.json | Description |
|-----|------|-------------------------------|-------------|
| `network` | string | `files/networks/UKNet_BDM.json` | Path to network topology |
| `routes` | string | `files/routes/UKNet_routes.json` | Path to pre-calculated routes |
| `capacities` | string | `files/capacities/capacities.json` | Path to band capacities |
| `bitrate` | string | `files/bitrate/bitrate.json` | Path to bitrate configuration |
| `lambda` | number | `600` | Arrival rate $\lambda$ |
| `mu` | number | `1` | Service rate $\mu$ |
| `bands` | integer | `4` | Number of frequency bands |
| `goal` | number | `1e7` | Number of connections to simulate |
| `logs` | boolean | `true` | Enable progress logging |
| `legacy` | boolean | `false` | Use legacy file loaders |
| `defrag_mode` | string | `none` | Defragmentation mode (`none`, `before_arrival`, `after_block`, `after_assign`) |
| `events_csv` | string | empty | Optional path to save generated event CSV |
| `sweep` | boolean | `true` | Run a lambda sweep instead of one simulation |
| `lambda_start` | number | `500` | Start lambda for a sweep |
| `lambda_end` | number | `1000` | End lambda for a sweep |
| `lambda_step` | number | `100` | Lambda increment for a sweep; must be greater than zero |

Notes:
- CLI flag `-logs` overrides the `logs` value from config.
- CLI flag `-events-csv` overrides `events_csv` from config.
- A sweep runs when either `sweep` in the JSON configuration is `true` or `-sweep=true` is provided.
- The sweep range is read from `lambda_start`, `lambda_end`, and `lambda_step` in the JSON configuration when those values are present; the CLI range flags provide fallback values.
- Because `sweep` is enabled when either source is `true`, `-sweep=false` does not disable a sweep enabled by the JSON configuration.

## Testing

To run the unit tests and verify the implementation:
```bash
make test
```

## Makefile Targets

- `make test`: run all tests (`go test -v ./...`).
- `make build`: clean and build `bin/simulador`.
- `make run`: build and execute `bin/simulador` with all CLI flag variables. Supported variables are `CONFIG`, `LOGS`, `EVENTS_CSV`, `SWEEP`, `LAMBDA_START`, `LAMBDA_END`, and `LAMBDA_STEP`.
- `make clean`: remove build output under `bin/`.

## Project Structure

- `main.go`: Entry point, handles CLI flags and starts the simulation.
- `internal/`: Core simulation logic and models.
    - `allocator/`: Resource allocation algorithms (e.g., FirstFit).
    - `infrastructure/`: Network, Nodes, Links, and Spectrum management.
    - `simulator/`: Main simulation engine.
- `pkg/`: Utility packages for plotting and validation.
- `files/`: Input configuration files in JSON format.
- `legacy_files/`: Legacy file configurations.
    - `bitrates/`: Legacy bitrate configurations.
    - `networks/`: Legacy network configurations and capacity.
    - `routes/`: Legacy route configurations.
- `result/`: Directory where generated plots and results are saved.
- `bin/`: Compiled binaries.

## Results

After a simulation finishes, a plot is automatically generated in the `result/` directory, showing the blocking probability vs. the number of connections.