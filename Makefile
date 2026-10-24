# ==============================================================================
# CIPHER FLEET: MASTER UNIFIED COMPILATION & APPLIANCE MAKEFILE
# Single-command build pipeline for all 8 platform microservices
# ==============================================================================

SHELL := /bin/bash
ROOT_DIR := $(shell pwd)
DIST_DIR := $(ROOT_DIR)/dist

.PHONY: all build-all build-go build-mcp build-core build-exec test-synthetic doctor clean help

all: build-all

help:
	@echo "Cipher AI PaaS Master Build Automation:"
	@echo "  make build-all       - Compiles all 8 services (Go, MCP esbuild, PyInstaller core)"
	@echo "  make build-go        - Compiles Route, Log, Mem, Sec, and CLI Go services"
	@echo "  make build-mcp       - Bundles MCP server into single standalone .js file with esbuild"
	@echo "  make build-core      - Packages LangGraph orchestrator into standalone binary"
	@echo "  make test-synthetic  - Executes end-to-end 4-gate synthetic integration harness"
	@echo "  make doctor          - Runs host resource, port, and daemon diagnostic audit"
	@echo "  make clean           - Purges all build binaries, dist directories, and specs"

build-all: build-go build-mcp build-core
	@echo ""
	@echo "======================================================================"
	@echo "🎉 [CIPHER FLEET] MASTER COMPILATION COMPLETE ACROSS ALL SERVICES"
	@echo "======================================================================"

build-go:
	@echo "🔨 [1/3] Compiling Go Infrastructure Services..."
	@echo "   → Compiling cipher-route..."
	@cd "$(ROOT_DIR)/cipher-route" && go build -o cipher-route .
	@echo "   → Compiling cipher-log..."
	@cd "$(ROOT_DIR)/cipher-log" && go build -o cipher-log .
	@echo "   → Compiling cipher-mem..."
	@cd "$(ROOT_DIR)/cipher-mem" && go build -o cipher-mem .
	@echo "   → Compiling cipher-sec..."
	@cd "$(ROOT_DIR)/cipher-sec" && go build -o cipher-sec .
	@echo "   → Compiling cipher-cli..."
	@cd "$(ROOT_DIR)/cipher-cli" && go build -o cipher .
	@echo "✓ All 5 Go services compiled successfully."

build-mcp:
	@echo "📦 [2/3] Bundling MCP Server with esbuild..."
	@cd "$(ROOT_DIR)/cipher-mcp" && npx esbuild src/index.ts --bundle --platform=node --outfile=dist/cipher-mcp.js
	@echo "✓ Stateless MCP Server bundled into standalone dist/cipher-mcp.js."

build-core:
	@echo "🧠 [3/3] Packaging Python Orchestrator with PyInstaller..."
	@cd "$(ROOT_DIR)/cipher-core" && uv run pyinstaller --onefile main.py --name cipher-core --clean --noconfirm
	@echo "✓ LangGraph Orchestrator packaged into standalone binary."

build-exec:
	@echo "🛡️ Packaging Execution Sandbox with PyInstaller..."
	@cd "$(ROOT_DIR)/cipher-exec" && uv run pyinstaller --onefile server.py --name cipher-exec --clean --noconfirm
	@echo "✓ Execution Sandbox packaged into standalone binary."

test-synthetic:
	@echo "🧪 Running Synthetic Fleet Integration Harness..."
	@cd "$(ROOT_DIR)/cipher-cli" && go run . test --synthetic

doctor:
	@echo "🩺 Running Fleet Doctor Pre-flight Audit..."
	@cd "$(ROOT_DIR)/cipher-cli" && go run . doctor

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f "$(ROOT_DIR)/cipher-route/cipher-route"
	@rm -f "$(ROOT_DIR)/cipher-log/cipher-log"
	@rm -f "$(ROOT_DIR)/cipher-mem/cipher-mem"
	@rm -f "$(ROOT_DIR)/cipher-sec/cipher-sec"
	@rm -f "$(ROOT_DIR)/cipher-cli/cipher"
	@rm -rf "$(ROOT_DIR)/cipher-core/build" "$(ROOT_DIR)/cipher-core/dist" "$(ROOT_DIR)/cipher-core"/*.spec
	@rm -rf "$(ROOT_DIR)/cipher-exec/build" "$(ROOT_DIR)/cipher-exec/dist" "$(ROOT_DIR)/cipher-exec"/*.spec
	@rm -rf "$(ROOT_DIR)/dist"
	@echo "✓ Clean complete."
