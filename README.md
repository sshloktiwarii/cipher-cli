# Foundry CLI: The Fleet Commander & Control Plane

A unified command-line control plane written in Go with `spf13/cobra` that manages the entire lifecycle of the Foundry Local AI PaaS—networking, container orchestration, and multiplexed colored telemetry.

## Commands

### 1. `foundry up`
Programmatically generates the unified `docker-compose.yml` merging:
- `foundry-sentinel` (`:8079`)
- `foundry-router` (`:8080`)
- `foundry-crucible` (`:8081`)
- `foundry-pulse` (`:8082`)
- `foundry-graph` (`:8083`)
- `foundry-airgap` (`:8085`)
- `timescaledb` (`:5432`)
- `redis` (`:6379`)

Spins up all containers in detached mode connected through the internal bridge network `foundry_net`.

### 2. `foundry logs`
Attaches to all running containers and multiplexes their stdout/stderr streams into a single terminal window with dedicated ANSI color codes per microservice.

### 3. `foundry status`
Displays the port mappings, container runtime state, and HTTP health status across the entire fleet.

### 4. `foundry halt`
Gracefully stops and removes all containers and network bridges.

## Building and Installing
```bash
go mod tidy
go build -o foundry main.go
# Optional: install to PATH
sudo mv foundry /usr/local/bin/
```
