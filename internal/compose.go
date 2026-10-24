package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateComposeSpec produces the master docker-compose.yml content for the Foundry fleet
func GenerateComposeSpec() string {
	return `version: '3.8'

networks:
  foundry_net:
    driver: bridge
    name: foundry_net

services:
  timescaledb:
    image: timescale/timescaledb:latest-pg16
    container_name: foundry-timescaledb
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: telemetry
    ports:
      - "5432:5432"
    networks:
      - foundry_net
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: foundry-redis
    ports:
      - "6379:6379"
    networks:
      - foundry_net
    restart: unless-stopped

  foundry-sentinel:
    build:
      context: ../foundry-sentinel
    container_name: foundry-sentinel
    ports:
      - "8079:8079"
    environment:
      - PORT=8079
      - TARGET_URL=http://foundry-router:8080
    networks:
      - foundry_net
    depends_on:
      - foundry-router
    restart: unless-stopped

  foundry-router:
    build:
      context: ../foundry-router
    container_name: foundry-router
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - REDIS_ADDR=redis:6379
      - TELEMETRY_ENDPOINT=http://foundry-pulse:8082/ingest
    networks:
      - foundry_net
    depends_on:
      - redis
      - foundry-pulse
    restart: unless-stopped

  foundry-crucible:
    build:
      context: ../foundry-crucible
    container_name: foundry-crucible
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - TELEMETRY_ENDPOINT=http://foundry-pulse:8082/ingest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - foundry_net
    depends_on:
      - foundry-pulse
    restart: unless-stopped

  foundry-pulse:
    build:
      context: ../foundry-pulse
    container_name: foundry-pulse
    ports:
      - "8082:8082"
    environment:
      - PORT=8082
      - DATABASE_URL=postgres://postgres:password@timescaledb:5432/telemetry?sslmode=disable
    networks:
      - foundry_net
    depends_on:
      - timescaledb
    restart: unless-stopped

  foundry-graph:
    build:
      context: ../foundry-graph
    container_name: foundry-graph
    ports:
      - "8083:8083"
    environment:
      - PORT=8083
      - VAULT_PATH=/vault
    volumes:
      - ../foundry-graph/sample_vault:/vault
    networks:
      - foundry_net
    restart: unless-stopped

  foundry-airgap:
    build:
      context: ../foundry-airgap
    container_name: foundry-airgap
    ports:
      - "8085:8085"
    environment:
      - HITL_PORT=8085
    networks:
      - foundry_net
    restart: unless-stopped
`
}

// WriteComposeFile persists the compose file to disk
func WriteComposeFile(targetDir string) (string, error) {
	composePath := filepath.Join(targetDir, "docker-compose.yml")
	content := GenerateComposeSpec()
	if err := os.WriteFile(composePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}
	return composePath, nil
}
