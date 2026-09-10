# AGENTS_README.md — AI-Agent Codebase Guide

Repository root: `/home/erick/tesis/simulador`
Module path: `github.com/Kayres21/optical-mb-sim-go`

This file is a structured reference for AI coding agents. It documents architecture, module boundaries, data flow, file formats, and conventions. All paths are relative to the repository root unless stated otherwise.

---

## 1. SYSTEM OVERVIEW & ARCHITECTURE

### 1.1 Core Purpose
A discrete-event simulator that evaluates **blocking probability** in optical multiband (C/L/S/E-band) elastic optical networks (EONs). It is a Go port of the "Flex Net Sim" C++/Python project. Given a network topology, precomputed routes, per-band spectrum capacities, and a bitrate/modulation catalogue, it simulates connection arrivals/departures (M/M/∞-style Poisson process), attempts spectrum allocation (RMSA — Routing, Modulation, and Spectrum Assignment) via a pluggable allocator (default: First-Fit with distance-adaptive modulation), optionally runs defragmentation strategies, and reports the connection-blocking probability with confidence intervals (Wald, Agresti-Coull, Wilson), plus PNG plots of blocking probability vs. number of connections.

### 1.2 Tech Stack
- **Language**: Go 1.24.4 (single module, no external services).
- **Runtime dependencies** (from [go.mod](go.mod)):
  - `github.com/google/uuid` — UUID generation.
  - `github.com/xeipuuv/gojsonschema` (+ `gojsonpointer`, `gojsonreference`) — JSON Schema validation of input files.
  - `gonum.org/v1/plot` (+ `go-fonts/liberation`, `go-latex/latex`, `go-pdf/fpdf`, `ajstarks/svgo`, `golang/freetype`, `sbinet/gg`, `golang.org/x/image`, `golang.org/x/text`) — PNG chart generation.
- **No database.** All state is in-memory per simulation run; inputs/outputs are JSON/CSV/PNG files on disk.
- **No web server / no API.** This is a CLI batch simulation tool.
- **Build tooling**: standard `go build`/`go test`, orchestrated via [Makefile](Makefile).

### 1.3 Architecture Pattern
Layered / Clean-Architecture-inspired CLI application:

```
main.go (CLI + config)  →  simulations/ (batch/sweep orchestration)
                         →  internal/simulator (event-loop engine)
                             ├─ internal/connections/controller (allocation + state)
                             │   ├─ internal/allocator (RMSA algorithms, function type)
                             │   ├─ internal/connections (domain models: Connection, Routes, BitRate, Events)
                             │   └─ internal/infrastructure (Network, Node, Link, Capacity/Band, fragmentation)
                             ├─ internal/defragmentator (spectrum defragmentation strategies)
                             └─ internal/connections/randomVariable (seeded RNG streams)
                         →  internal/loader (ResourceLoader: Standard vs. Legacy file format adapters)
                         →  pkg/validator (JSON Schema validation)
                         →  pkg/helpers (blocking probability + confidence interval math)
                         →  pkg/plotter (gonum-based PNG chart rendering)
```

Key design patterns:
- **Strategy pattern**: `allocator.Allocator` is a function type (not an interface) — swap RMSA algorithms by passing a different function.
- **Strategy pattern**: `loader.ResourceLoader` interface — `StandardLoader` (schema-validated JSON) vs. `LegacyLoader` (legacy JSON, no schema).
- **Strategy pattern**: `defragmentator.DecisionFunc` / `ActionFunc` — pluggable "when to defrag" and "how to defrag" behavior, selected by `defrag_mode`.
- **Event-driven simulation core**: `internal/simulator` uses a min-heap (`container/heap`) of `ConnectionEvent` ordered by time — classic discrete-event simulation (DES).
- **Dependency injection via constructors**: `simulator.New` / `NewWithSeeds`, `controller.New` wire allocator/routes/network together explicitly; no global state or DI framework.

---

## 2. FOLDER STRUCTURE & NAVIGATION

```
/home/erick/tesis/simulador/
├── main.go                        # CLI entry point: flags, JSON config loading, wiring, single-run or sweep dispatch
├── Makefile                       # build/test/run/clean targets
├── go.mod / go.sum                # Go module + dependency lockfile
├── README.md                      # human-oriented usage docs (CLI flags, config schema, examples)
├── bin/                           # compiled binary output (bin/simulador) — build artifact, gitignore candidate
├── result/                        # generated PNG plots (output of a run) — created at runtime
│
├── internal/                      # non-exported application logic (cannot be imported by external modules)
│   ├── allocator/                 # RMSA algorithms
│   │   ├── allocator.go           # Allocator func type + FirstFit implementation (distance-adaptive modulation)
│   │   ├── file_allocator.go      # FirstFitFromFile: replay allocation decisions from a CSV (for testing/determinism)
│   │   └── allocator_test.go
│   ├── connections/                # domain models for traffic/connections
│   │   ├── connections.go         # Connection struct (the "allocated resource" record)
│   │   ├── connectionEvents.go     # ConnectionEvent struct + GenerateEvents (bulk arrival generator)
│   │   ├── bitrate.go              # BitRate / BitRateList models + standard-format JSON loader (files/bitrate/*)
│   │   ├── routes.go               # Routes / Path models + standard-format JSON loader (files/routes/*)
│   │   ├── bitrate_test.go, connections_test.go, routes_test.go
│   │   ├── controller/
│   │   │   └── controller.go       # Controller: owns Network+Routes+Allocator+Connections map; allocation/release API
│   │   └── randomVariable/
│   │       └── randomVariable.go   # RandomVariable: seeded exponential/uniform RNG streams (arrive, departure, bitrate, source, dest, band)
│   ├── defragmentator/
│   │   ├── defragmentator.go       # DecisionFunc/ActionFunc types, DefragMode constants, FirstFitActiveConnections strategy
│   │   └── defragmentator_test.go
│   ├── infrastructure/             # physical network model
│   │   ├── network.go              # Network struct, ReadNetworkFile, NetworkGenerate (network+capacities merge)
│   │   ├── node.go                 # Node struct ({ID int})
│   │   ├── link.go                 # Link struct: per-band slot occupancy (bool slices), Assign/ReleaseConnection, FR1 fragmentation ratio
│   │   ├── capacities.go           # Capacity/Band struct + ReadCapacityFile (per-band spectrum slot counts)
│   │   └── *_test.go
│   ├── loader/                     # file-format adapters (Adapter pattern)
│   │   ├── loader.go                # ResourceLoader interface
│   │   ├── standard.go              # StandardLoader: delegates to infrastructure/connections schema-validated readers
│   │   ├── legacy.go                # LegacyLoader: parses legacy_files/-style JSON (no schema validation)
│   │   └── legacy_test.go
│   └── simulator/
│       ├── simulator.go            # Simulator struct: DES engine, event heap, RNG init, blocking-probability table, CI math, CSV export, Plot
│       └── simulator_test.go
│
├── pkg/                             # exported/reusable utility packages (importable outside this module)
│   ├── pkg.go                       # (package doc / marker file)
│   ├── helpers/
│   │   ├── helpers.go               # ComputeBlockingProbabilities, Wald/AgrestiCoull/Wilson confidence intervals
│   │   └── helpers_test.go, helpers_ci_test.go
│   ├── plotter/
│   │   ├── scatter_plot.go          # GenerateScatterPlot / GenerateLinePlot / GenerateMultiSeriesPlot (gonum/plot wrappers)
│   │   └── scatter_plot_test.go
│   └── validator/
│       └── validator.go             # Validate / ValidateFile: JSON Schema validation via gojsonschema
│
├── simulations/
│   └── simulations.go               # SweepRunner: runs the simulator across a lambda range, aggregates + multi-series plots
│
├── examples/
│   └── defrag_example.go            # standalone example program demonstrating defragmentation usage
│
├── files/                           # STANDARD (current) input file format — schema-validated JSON
│   ├── config.json                  # default AppConfig used by `-config` flag
│   ├── config-legacy.json, config-defrag-before-arrival.json, config-standard-same-data.json
│   ├── bitrate/                     # bitrate.json (+ legacy/newformat/test variants), schema.json, general.json
│   ├── capacities/                  # capacities.json (+ legacy/test variants), schema.json, general.json
│   ├── networks/                    # ARPANet/EON/Eurocore/NSFNet/UKNet/USNet/network_test .json topologies, schema.json
│   └── routes/                      # precomputed k-shortest-path routes per network, schema.json
│
├── legacy_files/                    # LEGACY input file format (read only by loader.LegacyLoader, no schema)
│   ├── bitrates/                    # bitrate_bpsk.json, bitrate_iroBand_C.json, bitrate_mb.json, bitrate.json
│   ├── networks/                    # legacy topology JSON (per-link "slots" as int or map[band]int)
│   └── routes/                      # legacy route JSON
│
└── test/                            # additional integration-style test fixtures/programs
    ├── main.go
    ├── e001/, events/, examples/, results/, routes/, test_ff/
```

### 2.1 Where to look for X
| Need to change... | Look in |
|---|---|
| CLI flags / config keys | [main.go](main.go) |
| RMSA / spectrum allocation algorithm | [internal/allocator/allocator.go](internal/allocator/allocator.go) |
| Network/link/spectrum model, fragmentation ratio (FR1) | [internal/infrastructure/link.go](internal/infrastructure/link.go), [internal/infrastructure/network.go](internal/infrastructure/network.go) |
| Connection lifecycle / bookkeeping | [internal/connections/controller/controller.go](internal/connections/controller/controller.go) |
| Event loop, blocking probability computation, CI math wiring | [internal/simulator/simulator.go](internal/simulator/simulator.go) |
| Defragmentation strategies | [internal/defragmentator/defragmentator.go](internal/defragmentator/defragmentator.go) |
| Reading standard vs legacy JSON input files | [internal/loader/standard.go](internal/loader/standard.go), [internal/loader/legacy.go](internal/loader/legacy.go) |
| Bitrate/modulation catalogue | [internal/connections/bitrate.go](internal/connections/bitrate.go) |
| Routes/paths format | [internal/connections/routes.go](internal/connections/routes.go) |
| Confidence interval / blocking probability formulas | [pkg/helpers/helpers.go](pkg/helpers/helpers.go) |
| Plot generation | [pkg/plotter/scatter_plot.go](pkg/plotter/scatter_plot.go) |
| Lambda sweeps (multi-run experiments) | [simulations/simulations.go](simulations/simulations.go) |
| JSON schema validation | [pkg/validator/validator.go](pkg/validator/validator.go) |
| Input JSON schemas | `files/*/schema.json` |
| Sample network/route/bitrate/capacity data | `files/networks/`, `files/routes/`, `files/bitrate/`, `files/capacities/` |

---

## 3. ENTRY POINTS & DATA FLOW

### 3.1 Application Entry Point
[main.go](main.go) `func main()`:
1. Parses CLI flags (`-config`, `-events-csv`, `-logs`, `-defrag-mode`, `-sweep`, `-lambda-start/end/step`).
2. Loads JSON config file (default `files/config.json`) into `AppConfig`, applies defaults via `applyDefaults`.
3. Validates `defrag_mode` against `defragmentator.DefragNone/DefragBeforeArrival/DefragAfterBlock/DefragAfterAssign`.
4. Selects `loader.ResourceLoader`: `StandardLoader` (default) or `LegacyLoader` (if `legacy: true`).
5. Loads network (`LoadNetwork`), bitrate catalogue (`LoadBitRate`), and routes (`LoadRoutes`).
6. Branches:
   - **Sweep mode** (`sweep: true` in config OR `-sweep=true`): builds `simulations.SweepRunner`, runs multiple simulations across a lambda range, produces one combined multi-series plot.
   - **Single-run mode**: constructs `simulator.Simulator` via `simulator.New(...)`, calls `sim.Start(logs)`, optionally writes an events CSV, then generates a single PNG plot.

Other entry points:
- [test/main.go](test/main.go) — separate test/integration driver.
- [examples/defrag_example.go](examples/defrag_example.go) — standalone example demonstrating defragmentation.

### 3.2 Primary Data Flow (single simulation run)

```
JSON config file ─┐
CLI flags ─────────┼──► AppConfig (main.go) ──► ResourceLoader.Load{Network,BitRate,Routes}
                   │                                   │
                   │                    (validates against files/*/schema.json via pkg/validator)
                   │                                   ▼
                   │                    infrastructure.Network, connections.BitRateList, connections.Routes
                   │                                   │
                   └──► simulator.New(...) ─────────────┘
                              │
                              ▼
                    Simulator.Start(logOn bool)
                              │
              ┌───────────────┴────────────────────────────────────────┐
              │  Event loop (min-heap ordered by Time)                 │
              │  1. pop earliest ConnectionEvent                       │
              │  2. if Arrive:                                         │
              │     - increment totalConnections, print progress row   │
              │     - (optional) DefragBeforeArrival defrag pass       │
              │     - schedule next random Arrive event                │
              │     - Controller.ConnectionAllocation(...) → allocator.Allocator (FirstFit)│
              │     - if blocked & DefragAfterBlock: defrag, retry once│
              │     - if assigned: schedule a Release event, ++assignedConnections│
              │       (optional) DefragAfterAssign defrag pass         │
              │  3. if Release: Controller.ReleaseConnection(...) frees slots on all links│
              └──────────────────────────────────────────────────────────┘
                              │
                              ▼
       results[] (blocking probability samples) + arrives[] (x-axis sample points)
                              │
                              ▼
              pkg/helpers (Wald/Agresti-Coull/Wilson CI) for progress table
                              │
                              ▼
              pkg/plotter.GenerateScatterPlot → result/<title>_<timestamp>.png
```

Allocation sub-flow (`allocator.FirstFit`):
```
Controller.ConnectionAllocation(src, dst, bitRate, numberOfBands, id)
  → Allocator(src, dst, bitRate, Network, Routes, numberOfBands, id, addConnection callback)
     → Routes.GetPaths(src, dst)                     # candidate k-shortest paths
     → for each path: Network.GetLinkByPath(...)      # resolve *Link objects
     → BitRate.GetDistanceAdaptive(pathLength)         # pick modulation with smallest reach ≥ length
     → for each band 0..numberOfBands-1:
         → aggregate per-link slot occupancy (First-Fit contiguous-slot search across all links on the path)
         → on fit: Link.AssignConnection(initialSlot, slotCount, band) for every link on path
         → addConnection(Connection{...}) callback → Controller.AddConnection → stored in Controller.Connections map
     → returns true/false (assigned or blocked)
```

Release sub-flow: `Controller.ReleaseConnection` looks up the stored `Connection` by ID, calls `Link.ReleaseConnection` on every link it used, invokes an optional `UnassignCallback`, then deletes it from the `Connections` map.

### 3.3 Sweep / Multi-run Data Flow
[simulations/simulations.go](simulations/simulations.go) `SweepRunner.Run()`: loops `lambda` from `StartLambda` to `EndLambda` in `StepLambda` increments, constructs and runs a fresh `simulator.Simulator` per lambda value (defrag disabled, `allocator.FirstFit`), collects `RunResult{Lambda, Arrives, Results}`, then `Plot()` renders all series on one chart via `plotter.GenerateMultiSeriesPlot`.

---

## 4. CORE MODULES & STATE MANAGEMENT

### 4.1 Business Domains
| Module | Responsibility |
|---|---|
| `internal/infrastructure` | Physical topology: `Network` (Nodes+Links), `Link` (per-band boolean slot occupancy + FR1 fragmentation ratio), `Capacity`/`Band` (spectrum slot counts per band). |
| `internal/connections` | Traffic domain: `Connection` (an allocated resource: links, slots, band), `ConnectionEvent` (Arrive/Release DES events), `BitRate`/`BitRateList` (modulation/reach/slots catalogue), `Routes`/`Path` (precomputed k-shortest paths). |
| `internal/connections/controller` | Orchestrates allocation: holds `Network`, `Routes`, `Allocator`, and the live `Connections map[string]Connection` (the simulation's mutable state). |
| `internal/connections/randomVariable` | Seeded RNG streams (`math/rand`) for arrival/departure (exponential) and bitrate/source/destination/band (uniform) sampling — one `*rand.Rand` per stream for reproducibility. |
| `internal/allocator` | RMSA strategy functions: `FirstFit` (live simulation) and `FirstFitFromFile` (CSV-replay for deterministic testing). |
| `internal/defragmentator` | Spectrum defragmentation strategies: decision (`when to trigger`) + action (`how to reallocate`) function pairs. |
| `internal/loader` | Adapts two on-disk file formats (standard schema-validated vs. legacy) into the same domain models. |
| `internal/simulator` | The DES engine: event heap, RNG orchestration, blocking-probability sampling/reporting, confidence intervals, CSV/plot export. |
| `simulations` | Batch orchestration: lambda sweeps across many simulator runs. |
| `pkg/helpers`, `pkg/plotter`, `pkg/validator` | Cross-cutting utilities (math, charting, schema validation). |

### 4.2 State Management & Persistence
- **No database.** All simulation state is in-memory, held by value/pointer inside `Simulator` and `Controller` structs for the lifetime of a single `main()` invocation (or one iteration of a sweep).
- **Mutable state root**: `Controller.Connections map[string]connections.Connection` — the authoritative record of active (allocated) connections; mutated by `AddConnection` (on allocation) and `delete()` (on release, inside `ReleaseConnection`).
- **Spectrum state**: each `infrastructure.Link` owns its own `Capacity.Bands[i].Slots []bool` (per-slot occupancy) plus `FragmentationRatioByBand []float64` (derived, recomputed on every `AssignConnection`/`ReleaseConnection` call — see `Link.UpdateFragmentationRatio`).
- **RNG state**: each `RandomVariable` stream holds its own `*rand.Rand` seeded independently (`SetSeeds`) — enables reproducible runs; default seeds are hardcoded constants in [internal/simulator/simulator.go](internal/simulator/simulator.go) (`defaultSeedArrive` etc.).
- **Persistence to disk** happens only at the boundaries:
  - **Input**: JSON files under `files/` or `legacy_files/` (read-only, loaded once at startup).
  - **Output**: PNG plots written to `result/` (`GenerateLinePlot`/`GenerateMultiSeriesPlot`), and optionally a CSV event log written to the path given by `events_csv` config key / `-events-csv` flag (`Simulator.SaveEventsCSV`).

---

## 5. API & INTERFACE SPECIFICATIONS

This is a CLI batch tool with **no network API, no webhooks, no REST/gRPC endpoints**. The public "interfaces" are:

### 5.1 CLI Flags (see [main.go](main.go))
| Flag | Default | Description |
|---|---|---|
| `-config` | `files/config.json` | Path to JSON configuration file |
| `-events-csv` | `""` | Path to write generated events CSV (overrides `events_csv` config key) |
| `-logs` | `true` | Enable progress logging |
| `-defrag-mode` | `""` | `none` \| `before_arrival` \| `after_block` \| `after_assign` (overrides config) |
| `-sweep` | `false` | Run a lambda sweep instead of a single simulation |
| `-lambda-start` / `-lambda-end` / `-lambda-step` | `500` / `1500` / `50` | Sweep range fallback values (JSON config values take precedence if present) |

### 5.2 JSON Configuration Schema (`AppConfig` in main.go)
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
  "events_csv": "",
  "sweep": false,
  "lambda_start": 500,
  "lambda_end": 1000,
  "lambda_step": 100
}
```
All keys are optional; missing values fall back to `defaultConfig()` in [main.go](main.go).

### 5.3 Internal Go "APIs" (function-type strategy interfaces)
- **Allocator** (`internal/allocator/allocator.go`):
  ```go
  type Allocator func(source, destination int, bitRate connections.BitRate, network infrastructure.Network,
      path connections.Routes, numberOfBands int, id string, addConnection func(connections.Connection)) bool
  ```
- **ResourceLoader** (`internal/loader/loader.go`):
  ```go
  type ResourceLoader interface {
      LoadNetwork(networkPath, capacitiesPath string) (infrastructure.Network, error)
      LoadBitRate(bitRatePath string, numberOfBands int) (connections.BitRateList, error)
      LoadRoutes(routesPath string) (connections.Routes, error)
  }
  ```
- **Defragmentator** (`internal/defragmentator/defragmentator.go`):
  ```go
  type DecisionFunc func(network infrastructure.Network, connections map[string]connections.Connection, event connections.ConnectionEvent, numberOfBands int) bool
  type ActionFunc func(network infrastructure.Network, connections map[string]connections.Connection, routes connections.Routes, alloc allocator.Allocator, numberOfBands int) error
  ```

### 5.4 Output CSV Schema (events, `-events-csv`)
Header: `ID,Tiempo,Evento,Source,Destination,BitRate` — `Evento` is `ARRIVE` or `DEPARTURE`.

### 5.5 Output PNG
Filename pattern: `<title>_<YYYYMMDD_HHMMSS>.png`, saved to `result/`. Titles are generated as `FirstFit_<NetworkAlias>-erlang-<lambda>_<bands>` (single run) or `FirstFit_<NetworkAlias>_lambda_sweep_<start>_<end>` (sweep).

---

## 6. DATA MODELS & SCHEMA

There is no database; "schema" here means JSON file formats validated by `pkg/validator` against `files/*/schema.json` (standard format only — `legacy_files/` has no schema and is parsed leniently by `LegacyLoader`).

### 6.1 Network file (`files/networks/*.json`, schema at `files/networks/schema.json`)
```json
{
  "name": "...", "alias": "...",
  "nodes": [{ "id": 0 }, ...],
  "links": [{ "id": 0, "src": 0, "dst": 1, "length": 100 }, ...]
}
```
Maps to `infrastructure.Network{Name, Alias, Nodes []Node, Links []Link}`. Per-link spectrum `Capacity` is injected separately from the capacities file via `NetworkGenerate` (not stored in the network JSON itself).

### 6.2 Capacities file (`files/capacities/*.json`)
```json
{ "bands": [{ "id": "0", "name": "C", "capacity": 320 }, ...] }
```
Maps to `infrastructure.Capacity{Bands []Band}`; `Band.Slots []bool` (size `SlotsLen`/`capacity`) is allocated at load time — this is the per-link, per-band spectrum grid, cloned onto every `Link` in `NetworkGenerate`.

### 6.3 Routes file (`files/routes/*.json`)
```json
{ "alias": "...", "name": "...", "routes": [{ "src": 0, "dst": 1, "paths": [[0,2,1], [0,3,1]] }] }
```
Maps to `connections.Routes{Alias, Name, Paths []Path}`; each `Path.PathLinks` is a list of alternative node-ID sequences (k-shortest paths) between `Source`/`Destination`.

### 6.4 Bitrate file (`files/bitrate/*.json`, schema-validated)
Grouped by modulation, each with `slots` (per Gb/s magnitude) and `reachs` (per number-of-bands configuration, with per-band reach breakdown). Converted by `convertSchemaBitRateFile` into one `connections.BitRate` per magnitude (Gbps value), each holding parallel slices: `Modulation[]`, `Slots[]`, `Reach[]`, and optional `Bands[][]`, `SlotsPerBand[][]`, `ReachPerBand[][]` for multi-band files.

### 6.5 Legacy formats (`legacy_files/`, no schema validation)
- Network: `{"Name":..., "alias":..., "nodes":[{"id":...}], "links":[{"id","src","dst","length","slots"}]}` where `slots` is either a plain int (single band, always named `"C"`) or `map[string]int` (band name → slot count, sorted alphabetically for stable band ordering).
- Bitrate: `map[gigabits string][]map[modulation string]config` where `config` is either `{"slots","reach"}` (single band) or `[{"<band>": {"slots","reach"}}, ...]` (multi-band). Modulations are sorted by a fixed rank (BPSK < QPSK < 8QAM < 16QAM < others).

### 6.6 In-memory domain entities & relationships
```
Network 1─* Node
Network 1─* Link ── Capacity 1─* Band (Slots []bool, FragmentationRatioByBand []float64)
Routes 1─* Path (Source, Destination, PathLinks [][]int)
BitRateList 1─* BitRate (Modulation[], Slots[], Reach[], optional per-band slices)
Connection ── *Link (many, one per hop) ; Source, Destination, InitialSlot, FinalSlot, Slots, BandSelected
Controller ── Network, Routes, Allocator, Connections map[string]Connection
ConnectionEvent (Arrive|Release) ── Id, Source, Destination, Bitrate index, Time, optional Connection* pointer fields
```

---

## 7. DEVELOPMENT, TESTING & DEPLOYMENT RUNBOOKS

### 7.1 Prerequisites
- Go ≥ 1.24 (see [go.mod](go.mod)).
- `make` (optional, recommended).

### 7.2 Install dependencies
```bash
go mod download
```

### 7.3 Run locally
```bash
make run                                                   # default config: files/config.json
make run CONFIG=files/config-legacy.json
make run CONFIG=files/config.json LOGS=false EVENTS_CSV=files/sim_events.csv
make run SWEEP=true LAMBDA_START=500 LAMBDA_END=1500 LAMBDA_STEP=50
# or directly:
go run main.go -config files/config.json
go run main.go -config files/config.json -sweep=true -lambda-start 500 -lambda-end 1500 -lambda-step 50
```

### 7.4 Run tests
```bash
make test          # equivalent to: go test -v ./...
```
Test files exist alongside almost every package (`*_test.go`) — unit tests cover `allocator`, `connections` (bitrate/routes/connections), `defragmentator`, `infrastructure` (capacities/link/network/node), `simulator`, `pkg/helpers`, `pkg/plotter`. There is no separate integration-test command beyond `go test ./...`; `test/` and `test_ff/` contain standalone fixture programs, not `go test` suites.

### 7.5 Build
```bash
make build          # cleans bin/, then: go build -o bin/simulador .
```

### 7.6 Clean
```bash
make clean          # rm -rf bin/
```

### 7.7 Environment variables
None required by the application itself — all configuration is via the `-config` JSON file and CLI flags. No secrets, API keys, or credentials are used anywhere in this codebase.

### 7.8 Deployment
This is a CLI simulation tool, not a deployed service. "Deployment" = distributing the compiled `bin/simulador` binary (or running via `go run`) alongside the `files/`/`legacy_files/` input data directories, which must remain accessible at the paths referenced by the config JSON (paths are relative to the process's working directory, not embedded in the binary).

---

## 8. DESIGN PATTERNS, CONVENTIONS & CONSTRAINTS

- **Package layout convention**: `internal/` for application-private logic (enforced by Go's `internal/` visibility rule — cannot be imported by other modules), `pkg/` for reusable, potentially externally-importable utilities.
- **Strategy-via-function-type**, not interfaces, for `Allocator`, `DecisionFunc`, `ActionFunc` — prefer this pattern when adding new pluggable algorithms in this codebase, matching existing conventions.
- **Reproducibility constraint**: RNG streams (`randomVariable.RandomVariable`) are always explicitly seeded (`SetSeeds`); default seeds are named constants in [internal/simulator/simulator.go](internal/simulator/simulator.go) (`defaultSeedArrive` = 1, etc.) — never use unseeded/global `math/rand` in simulation-affecting code paths, to preserve deterministic, comparable simulation runs.
- **Band count constraint**: `numberOfBands` is clamped to `[minBands=1, maxBands=4]` in `Simulator.initVariableNumbers` — out-of-range values are logged as warnings and clamped, not rejected as errors.
- **Error handling pattern**: constructors and file-loading functions return `(T, error)` and wrap errors with `fmt.Errorf("...: %w", err)` for context; the CLI (`main.go`) uses `log.Fatalf` on unrecoverable startup errors (fail-fast at the boundary). Inside the simulation loop, per-event failures (e.g., failed release) are logged via `slog.Warn`/`slog.Error` and skipped rather than aborting the whole run.
- **Validation boundary**: all "standard" format JSON input files are schema-validated at the point of read (`pkg/validator.ValidateFile`) using a sibling `schema.json` in the same directory; legacy-format files are intentionally NOT schema-validated (`loader.LegacyLoader`) — do not add schema validation to legacy loaders without a corresponding schema file and explicit request.
- **Mutation via pointer slices**: `Connection.Links []*infrastructure.Link` and `Network.GetLinkByPath` return pointers so that `Link.AssignConnection`/`ReleaseConnection` mutate the single shared `Link` instances owned by `Network.Links` — do not copy `Link` by value when mutating spectrum state.
- **Fragmentation ratio (FR1)** is auto-maintained: any code path that mutates slot occupancy must go through `Link.AssignConnection`/`ReleaseConnection` (which call `UpdateFragmentationRatio` internally) rather than mutating `Capacities.Bands[i].Slots` directly, or the fragmentation ratio will become stale (see repo memory note `link-fragmentation.md`).
- **Allocator callback contract**: any `Allocator` implementation must invoke the `addConnection` callback on success so the `Controller` persists connection state — a `bool`-only return with no callback invocation will silently leave `Controller.Connections` inconsistent (see repo memory note `allocator-api.md`).
- **No concurrency**: the simulator is single-threaded/sequential (one event heap, no goroutines in the hot path) — do not introduce concurrent mutation of `Network`/`Link`/`Controller.Connections` without redesigning for synchronization.
- **Naming**: Spanish is used for some CLI output/log strings and axis labels (e.g., "Número de conexiones", "Probabilidad de bloqueo") reflecting the thesis/academic context of this project; keep this in mind when matching existing string conventions in user-facing output.
- **Testing convention**: table-style unit tests per package, named `Test<FunctionName>_<Scenario>` (e.g., `TestFirstFit_FailContiguity`); prefer constructing minimal in-memory `infrastructure.Network`/`connections.Routes` fixtures directly in test files rather than reading from `files/`.
