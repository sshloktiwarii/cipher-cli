package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateComposeSpec produces the master docker-compose.yml content for the Cipher fleet
func GenerateComposeSpec() string {
	return `version: '3.8'

networks:
  cipher_net:
    driver: bridge
    name: cipher_net

services:
  timescaledb:
    image: timescale/timescaledb:latest-pg16
    container_name: cipher-timescaledb
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: telemetry
    ports:
      - "5432:5432"
    networks:
      - cipher_net
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: cipher-redis
    ports:
      - "6379:6379"
    networks:
      - cipher_net
    restart: unless-stopped

  cipher-sec:
    build:
      context: ../cipher-sec
    container_name: cipher-sec
    ports:
      - "8079:8079"
    environment:
      - PORT=8079
      - TARGET_URL=http://cipher-route:8080
    networks:
      - cipher_net
    depends_on:
      - cipher-route
    restart: unless-stopped

  cipher-route:
    build:
      context: ../cipher-route
    container_name: cipher-route
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - REDIS_ADDR=redis:6379
      - TELEMETRY_ENDPOINT=http://cipher-log:8082/ingest
    networks:
      - cipher_net
    depends_on:
      - redis
      - cipher-log
    restart: unless-stopped

  cipher-exec:
    build:
      context: ../cipher-exec
    container_name: cipher-exec
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - TELEMETRY_ENDPOINT=http://cipher-log:8082/ingest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - cipher_net
    depends_on:
      - cipher-log
    restart: unless-stopped

  cipher-log:
    build:
      context: ../cipher-log
    container_name: cipher-log
    ports:
      - "8082:8082"
    environment:
      - PORT=8082
      - DATABASE_URL=postgres://postgres:password@timescaledb:5432/telemetry?sslmode=disable
    networks:
      - cipher_net
    depends_on:
      - timescaledb
    restart: unless-stopped

  cipher-mem:
    build:
      context: ../cipher-mem
    container_name: cipher-mem
    ports:
      - "8083:8083"
    environment:
      - PORT=8083
      - VAULT_PATH=/vault
    volumes:
      - ../cipher-mem/sample_vault:/vault
    networks:
      - cipher_net
    restart: unless-stopped

  cipher-mcp:
    build:
      context: ../cipher-mcp
    container_name: cipher-mcp
    ports:
      - "8085:8085"
    environment:
      - HITL_PORT=8085
    networks:
      - cipher_net
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
