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

### Using Go directly

You can run the binary directly with the supported flags:
```bash
go run main.go -config files/config.json
```

Optional flags:
```bash
go run main.go -config files/config.json -logs=false
go run main.go -config files/config.json -events-csv files/sim_events.csv
```

## Configuration Parameters

Configuration is provided via JSON. The default file is `files/config.json`:

```json
{
    "network": "files/networks/UKNet_BDM.json",
    "routes": "files/routes/UKNet_routes.json",
    "capacities": "files/capacities/capacities.json",
    "bitrate": "files/bitrate/bitrate.json",
    "lambda": 400,
    "mu": 1,
    "bands": 1,
    "goal": 1e5,
    "logs": true,
    "legacy": false,
    "defrag_mode": "none"
}
```

| Key | Type | Default in files/config.json | Description |
|-----|------|-------------------------------|-------------|
| `network` | string | `files/networks/UKNet_BDM.json` | Path to network topology |
| `routes` | string | `files/routes/UKNet_routes.json` | Path to pre-calculated routes |
| `capacities` | string | `files/capacities/capacities.json` | Path to band capacities |
| `bitrate` | string | `files/bitrate/bitrate.json` | Path to bitrate configuration |
| `lambda` | number | `400` | Arrival rate $\lambda$ |
| `mu` | number | `1` | Service rate $\mu$ |
| `bands` | integer | `1` | Number of frequency bands |
| `goal` | number | `1e5` | Number of connections to simulate |
| `logs` | boolean | `true` | Enable progress logging |
| `legacy` | boolean | `false` | Use legacy file loaders |
| `defrag_mode` | string | `none` | Defragmentation mode (`none`, `before_arrival`, `after_block`, `after_assign`) |
| `events_csv` | string | empty | Optional path to save generated event CSV |

Notes:
- CLI flag `-logs` overrides the `logs` value from config.
- CLI flag `-events-csv` overrides `events_csv` from config.

## Testing

To run the unit tests and verify the implementation:
```bash
make test
```

## Makefile Targets

- `make test`: run all tests (`go test -v ./...`).
- `make build`: clean and build `bin/simulador`.
- `make run`: build and execute `bin/simulador` with `-config="$(CONFIG)"`.
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