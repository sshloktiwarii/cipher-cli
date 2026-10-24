# cipher-cli: The Fleet Commander & Control Plane

A unified command-line control plane written in Go with `spf13/cobra` that manages the entire lifecycle of the Cipher Local AI PaaS—networking, container orchestration, and multiplexed colored telemetry.

## Commands

### 1. `cipher up`
Programmatically generates the unified `docker-compose.yml` merging:
- `cipher-sec` (`:8079`)
- `cipher-route` (`:8080`)
- `cipher-exec` (`:8081`)
- `cipher-log` (`:8082`)
- `cipher-mem` (`:8083`)
- `cipher-mcp` (`:8085`)
- `timescaledb` (`:5432`)
- `redis` (`:6379`)

Spins up all containers in detached mode connected through the internal bridge network `cipher_net`.

### 2. `cipher logs`
Attaches to all running containers and multiplexes their stdout/stderr streams into a single terminal window with dedicated ANSI color codes per microservice.

### 3. `cipher status`
Displays the port mappings, container runtime state, and HTTP health status across the entire fleet.

### 4. `cipher doctor`
Audits host resources, Docker daemon socket availability, port collisions across 8079-8085, and RAM availability for Apple Silicon unified memory or Linux hosts.

### 5. `cipher test --synthetic`
Executes end-to-end integration tests asserting the 4 core platform invariants:
- **Gate 1**: `cipher-sec` (:8079) credential intercept and Shannon entropy redaction.
- **Gate 2**: `cipher-exec` (:8081) Seatbelt process ceiling and runaway compute termination.
- **Gate 3**: `cipher-log` (:8082) SQLite WAL telemetry ingestion and event persistence.
- **Gate 4**: `cipher-mem` (:8083) sub-millisecond semantic search and Obsidian incident backlink resolution.

### 6. `cipher halt`
Gracefully halts and dismantles all running fleet containers and purges ephemeral network bridges.

## Zero-Dependency Distribution
Compile all Go microservices (`cipher-route`, `cipher-log`, `cipher-mem`, `cipher-cli`, `cipher-sec`) for `darwin-arm64` and `linux-amd64` using `CGO_ENABLED=0`, and package standalone Python executables (`cipher-core`, `cipher-exec`):
```bash
./scripts/build_distribution.sh
```

## Building the Binary Locally
```bash
go mod tidy
go build -o cipher .
./cipher status
./cipher test --synthetic
```

