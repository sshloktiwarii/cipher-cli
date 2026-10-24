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

### 4. `cipher halt`
Gracefully halts and dismantles all running fleet containers and purges ephemeral network bridges.

## Building the Binary
```bash
go mod tidy
go build -o cipher .
./cipher status
```
