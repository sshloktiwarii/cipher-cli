# 🛡️ Cipher AI PaaS — The Edge-Native Autonomous AI Operating System

[![Language](https://img.shields.io/badge/Language-Go%201.24-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Status](https://img.shields.io/badge/Status-Production%20Ready%20(10%2F10)-brightgreen?style=flat-square)](#)
[![Dependencies](https://img.shields.io/badge/Dependencies-Zero%20(No%20Docker%2C%20No%20Python%2C%20No%20Node)-blueviolet?style=flat-square)](#)
[![License](https://img.shields.io/badge/License-MIT-orange?style=flat-square)](LICENSE)

**Cipher** is a zero-dependency, edge-native autonomous AI platform. It acts as an automated safety cage, memory bank, and immune system for local AI agents (Ollama, Claude, Cursor, Antigravity)—blocking credential leaks, killing runaway compute loops, remembering incident resolutions, and proving tamper-resistance via cryptographic Merkle trees.

---

### 💡 For Non-Techies: What is Cipher in Plain English?

> Think of Cipher as an **airbag and flight recorder for AI running on your personal computer**.
>
> When you run modern AI tools locally, several risky things happen:
> - **They can leak secrets:** An AI might accidentally upload your private API keys, passwords, or personal documents to the internet.
> - **They can freeze your computer:** An AI writing code can trigger an infinite loop that consumes 100% of your computer's CPU and RAM until your system crashes.
> - **They forget everything:** AI models have no persistent memory—they don't remember how you solved a bug yesterday.
> - **Installation is painful:** Traditional AI platforms force you to install Python, Docker, Node.js, and gigabytes of complex packages.
>
> **How Cipher solves this:**
> 1. 🛡️ **Zero Leaks**: It spots and blacks out passwords in **1.4 milliseconds** before they can ever leave your machine.
> 2. ⏱️ **Automatic Kill-Switch**: If AI code gets stuck in an infinite loop, Cipher pulls the plug within **3 seconds**.
> 3. 🧠 **Permanent Private Memory**: Connects to your Obsidian notes folder to find relevant solutions in under **1 millisecond**.
> 4. 🚀 **Zero Installs**: The entire 8-service system compiles into a **single 35MB program**. You download it, run `cipher appliance`, and it boots in **30 milliseconds**. Zero Docker, zero Python, zero Node.js.

---

## ⚡ 1-Minute Quickstart

### Option A: Homebrew (macOS & Linux)
```bash
brew install sshloktiwarii/tap/cipher
```

### Option B: One-Line Installer
```bash
curl -fsSL https://raw.githubusercontent.com/sshloktiwarii/cipher-cli/master/scripts/install.sh | bash
```

### Option C: Build from Source
```bash
git clone https://github.com/sshloktiwarii/cipher-cli.git
cd cipher-cli
go build -o cipher .
```

---

## 🚀 Running the Cluster

### Start the Entire 8-Service Appliance (Single Supervised Command)
```bash
cipher appliance
```
*Boots the full cluster on Unix Domain Sockets in `< 500ms`, exposing only port `:8080`:*
```text
==================================================
🛡️  CIPHER EDGE APPLIANCE — IN-PROCESS SUPERVISOR
==================================================
✓ Initialized IPC Socket Cage at /tmp/cipher (0700 boundary)
✓ Started cipher-log     ➔ /tmp/cipher/log.sock (0600 perms)
✓ Started cipher-sec     ➔ /tmp/cipher/sec.sock (0600 perms)
✓ Started cipher-mem     ➔ /tmp/cipher/mem.sock (0600 perms)
✓ Started cipher-exec    ➔ /tmp/cipher/exec.sock (0600 perms)
✓ Started cipher-mcp     ➔ /tmp/cipher/mcp.sock (0600 perms)
✓ Started cipher-core    ➔ /tmp/cipher/core.sock (0600 perms)
✓ Started cipher-route   ➔ Ingress :8080

==================================================
🎉 CIPHER 10/10 APPLIANCE READY in 38ms
==================================================
• Unified Ingress:       http://localhost:8080 (Zero-Port Sprawl)
• War Room Cyber-HUD:    http://localhost:8080/hud
• Merkle Audit Proof:    http://localhost:8080/verify-audit-chain
• Cluster Topology:      http://localhost:8080/cluster/status
• Internal IPC:          unix:///tmp/cipher/*.sock
==================================================
```

### Run Synthetic Integration Verification
```bash
cipher test
```
*Verifies all 4 security and self-healing gates:*
```text
[GATE 1] Testing Semantic Firewall (cipher-sec :8079)...
   ✓ Gate 1 PASSED: Credential AKIA... intercepted & redacted in 1.46ms

[GATE 2] Testing Execution Sandbox Process Ceiling (cipher-exec :8081)...
   ✓ Gate 2 PASSED: Runaway compute killed via isolation (native-seatbelt) in 3.003s (ceiling <= 5s)

[GATE 3] Polling Neural Memory & Obsidian Vault (cipher-mem :8083)...
   ✓ Gate 3 PASSED: Verified incident post-mortem [[INC-*-resolution.md]] resolved in 801µs

[GATE 4] Verifying Cryptographic Merkle Hash Chain (cipher-log :8082 / :8080)...
   ✓ Gate 4 PASSED: Merkle chain verified untampered in 569µs

🎉 SYNTHETIC FLEET INTEGRATION TEST PASSED (All 4 verification gates satisfied)
```

### Self-Update from GitHub Releases
```bash
cipher update
```

---

## 🏛️ System Architecture

```text
                               +-----------------------------+
                               |     Client / AI Agent       |
                               | (Claude, Cursor, CLI, HUD)  |
                               +--------------+--------------+
                                              |
                                HTTP :8080 (Bearer Token)
                                              v
                      +-----------------------------------------------+
                      |                  cipher-route                 |
                      |   • Ingress Reverse Proxy & Port Aggregator   |
                      |   • Semantic Vector Cache (0.5ms Response)    |
                      |   • Real Cloud Circuit Breaker (Groq/OpenAI)  |
                      |   • Zero-Trust Session Secret (~/.cipher)     |
                      +-------+-----------+-----------+-------+-------+
                              |           |           |       |
            +-----------------+           |           |       +-----------------+
            |                             |           |                         |
unix:///tmp/cipher/sec.sock   unix:///tmp/cipher/exec.sock  unix:///tmp/cipher/mem.sock  unix:///tmp/cipher/log.sock
            |                             |           |                         |
            v                             v           v                         v
     +--------------+              +--------------+  +--------------+    +--------------+
     |  cipher-sec  |              | cipher-exec  |  |  cipher-mem  |    |  cipher-log  |
     | • WAF Filter |              | • Seatbelt   |  | • Vector KB  |    | • Merkle Tree|
     | • Redaction  |              | • 3s Timeout |  | • Obsidian   |    | • SQLite WAL |
     +--------------+              +--------------+  +--------------+    +--------------+
                                          ^                 ^
                                          |                 |
                             +------------+-----------------+------------+
                             |                cipher-core                |
                             |  • 5-Node Autonomous Remediation Machine  |
                             |  • Telemetry -> Retrieve -> Fix -> Test   |
                             +--------------------+----------------------+
                                                  |
                                                  v
                                     unix:///tmp/cipher/mcp.sock
                                     +---------------------------+
                                     |        cipher-mcp         |
                                     |  • Native Go MCP Server   |
                                     |  • JSON-RPC 2.0 Tools     |
                                     +---------------------------+
```

---

## 📦 The 8 Pure-Go Microservices

All 8 repositories are public, standalone, and compile with standard `go build .`:

| Service | Repository Link | Description | Stack |
| :--- | :--- | :--- | :--- |
| **`cipher-cli`** | [`sshloktiwarii/cipher-cli`](https://github.com/sshloktiwarii/cipher-cli) | Master supervisor, CLI control plane, self-updater & Homebrew formula. | Go, Cobra |
| **`cipher-route`** | [`sshloktiwarii/cipher-route`](https://github.com/sshloktiwarii/cipher-route) | Unified ingress gateway, semantic cache, circuit breaker & Cockpit HUD. | Go, HTTP, Embedded HUD |
| **`cipher-exec`** | [`sshloktiwarii/cipher-exec`](https://github.com/sshloktiwarii/cipher-exec) | Kernel-level ephemeral sandbox (Darwin Seatbelt & Linux namespaces). | Go, Seatbelt |
| **`cipher-mcp`** | [`sshloktiwarii/cipher-mcp`](https://github.com/sshloktiwarii/cipher-mcp) | Model Context Protocol server exposing isolated tools to LLMs. | Go, JSON-RPC 2.0, SQLite |
| **`cipher-core`** | [`sshloktiwarii/cipher-core`](https://github.com/sshloktiwarii/cipher-core) | 5-node autonomous self-healing state machine. | Go State Machine |
| **`cipher-sec`** | [`sshloktiwarii/cipher-sec`](https://github.com/sshloktiwarii/cipher-sec) | Semantic firewall, WAF regex & Shannon entropy secret interceptor. | Go, Regex Engine |
| **`cipher-mem`** | [`sshloktiwarii/cipher-mem`](https://github.com/sshloktiwarii/cipher-mem) | Obsidian vault watcher, vector embedding store & graph visualizer. | Go, fsnotify, Cosine Sim |
| **`cipher-log`** | [`sshloktiwarii/cipher-log`](https://github.com/sshloktiwarii/cipher-log) | Cryptographic append-only SHA-256 Merkle tree audit ledger. | Go, SQLite WAL, Merkle |

---

## 🔒 Zero-Trust Security Invariants

1. **Local Authentication Boundary**: On first boot, Cipher generates a 256-bit cryptographic secret in `~/.cipher/token` (`0600` permissions). All external HTTP calls to `/v1/chat/completions` or internal proxy routes require `Authorization: Bearer <token>`.
2. **Unix Domain Sockets (`0600`)**: Internal IPC takes place over `/tmp/cipher/*.sock`. File permissions are strictly pinned to `srw-------`, preventing other unprivileged local host processes from interacting with the backend cluster.
3. **Zero Hardcoded Mocks**: All circuit breakers fall back to real cloud providers (Groq/OpenAI) or return structured HTTP 503 errors when local models are down.

---

## 🛠️ CLI Commands Reference

```bash
cipher appliance          # Boots all 8 services in a single supervised process
cipher test --synthetic   # Runs the 4-gate integration verification test
cipher update             # In-place self update from GitHub Releases
cipher doctor             # Pre-flight host resource and socket audit
cipher status             # Displays health and port status across all services
cipher halt               # Gracefully reaps supervised processes and cleans sockets
```

---

## 📄 License

MIT License © 2026 [Shlok Tiwari](https://github.com/sshloktiwarii). Contributions welcome!
