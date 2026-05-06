# Optical Multiband Simulator

This simulator is a performance evaluation tool for optical multiband networks, based on the "Flex Net Sim" project. It evaluates network blocking probability under different arrival rates ($\lambda$), service rates ($\mu$), and frequency band configurations.

## Prerequisites

- **Go**: Version 1.24 or higher is required.
- **Make**: (Optional) For using the provided automation commands.

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd simulador
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Simulator

The recommended way to run the simulator is using the `Makefile`.

### Using Make (Recommended)

Run the simulation with default parameters:
```bash
make run
```

You can override any parameter directly from the command line:
```bash
make run BANDS=4 LAMBDA=100 GOAL=1e6
```

### Using Go directly

If you prefer to use `go run`, you can pass the flags manually:
```bash
go run main.go -bands 4 -lambda 100 -goal 1000000
```

## Configuration Parameters

| Parameter | Makefile Variable | Flag | Default | Description |
|-----------|-------------------|------|---------|-------------|
| Network | `NETWORK` | `-network` | `files/networks/UKNet_BDM.json` | Path to network topology |
| Routes | `ROUTES` | `-routes` | `files/routes/UKNet_routes.json` | Path to pre-calculated routes |
| Capacities| `CAPACITIES` | `-capacities` | `files/capacities/capacities.json`| Path to band capacities |
| Bitrate | `BITRATE` | `-bitrate` | `files/bitrate/bitrate.json` | Path to bitrate configurations|
| Lambda | `LAMBDA` | `-lambda` | `50` | Arrival rate $\lambda$ |
| Mu | `MU` | `-mu` | `1` | Service rate $\mu$ |
| Bands | `BANDS` | `-bands` | `1` | Number of frequency bands (1–4) |
| Goal | `GOAL` | `-goal` | `1e8` | Number of connections to simulate |
| Logs | `LOGS` | `-logs` | `true` | Enable real-time progress logging |

## Testing

To run the unit tests and verify the implementation:
```bash
make test
```

## Project Structure

- `main.go`: Entry point, handles CLI flags and starts the simulation.
- `internal/`: Core simulation logic and models.
    - `allocator/`: Resource allocation algorithms (e.g., FirstFit).
    - `infrastructure/`: Network, Nodes, Links, and Spectrum management.
    - `simulator/`: Main simulation engine.
- `pkg/`: Utility packages for plotting and validation.
- `files/`: Input configuration files in JSON format.
- `result/`: Directory where generated plots and results are saved.
- `bin/`: Compiled binaries.

## Results

After a simulation finishes, a plot is automatically generated in the `result/` directory, showing the blocking probability vs. the number of connections.