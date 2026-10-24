# Cipher: The Plain-English Guide

If you use AI tools (like Claude, Cursor, Ollama, or coding assistants) directly on your laptop, you run into three problems almost immediately:

1. **AI models can accidentally leak your secrets.** If an AI script reads a configuration file containing an API key or password, it might send that key over the internet to a cloud server without you realizing it.
2. **AI-generated code can freeze your machine.** An AI might write a script that enters an infinite loop, hogging 100% of your computer's CPU and RAM until your screen freezes.
3. **AI models forget everything.** When you close a chat, the AI forgets the bugs you already solved or the notes you took in your documentation.
4. **Setting up developer tools is frustrating.** Most local AI projects force you to install Python, Docker, Node.js, and hundreds of dependencies that break on different computers.

Cipher is a single program that solves all four issues on your machine.

---

## What Cipher Does For You

### 1. It blacks out passwords before they leave your computer
Cipher sits between your AI tool and the outside world. It scans text in about 1 millisecond. If it spots an AWS token, an API key, or a secret password, it replaces the key with `[REDACTED_BY_CIPHER_SEC]`. Your secrets stay on your machine.

### 2. It pulls the plug on runaway scripts
When an AI generates code and runs it, Cipher runs that code inside an isolated box with a 3-second hard timer. If the code gets stuck in an infinite loop, Cipher kills the process instantly. Your computer never freezes.

### 3. It gives your AI a permanent memory library
If you keep notes in Obsidian or Markdown files, Cipher indexes them locally. When a problem occurs, Cipher searches your notes by meaning (not just exact keywords) and hands the relevant past solution to the AI in less than 1 millisecond.

### 4. Zero installation headaches
Cipher is written in Go and compiles into a single file. You do not need Docker running in the background. You do not need Python or Node.js installed. You download it, double-click or run one command, and it boots in less than 40 milliseconds.

---

## How to Use It in 30 Seconds

### Step 1: Install
Open your terminal and run:

```bash
curl -fsSL https://raw.githubusercontent.com/sshloktiwarii/cipher-cli/master/scripts/install.sh | bash
```

*(On macOS or Linux with Homebrew, you can also run `brew install sshloktiwarii/tap/cipher`)*

### Step 2: Start the System
Type:

```bash
cipher appliance
```

You will see confirmation that the services started. Cipher is now running quietly in the background.

### Step 3: Open the Visual Dashboard
Open your web browser and go to:

```text
http://localhost:8080/hud
```

This opens a visual dashboard showing active services, memory searches, and security alerts in real time.

### Step 4: Stop it anytime
When you want to stop Cipher, press `Ctrl+C` in your terminal or run:

```bash
cipher halt
```

All background processes stop immediately and clean up after themselves.

---

## Frequently Asked Questions

#### Does Cipher upload my private files to the cloud?
No. All memory indexing, security scans, and code sandboxes run locally on your processor. Nothing is sent to any external server unless you explicitly configure a cloud model key.

#### Will this slow down my computer?
No. Cipher is written in native compiled Go, not interpreted Python or Electron. It uses less than 50MB of RAM and stays idle until an AI request comes through.

#### Do I need to be a programmer to use this?
No. If you can copy and paste `cipher appliance` into your terminal, you can run Cipher. The dashboard runs in your browser and displays what is happening in visual cards.

#### How do I know if Cipher is working?
Run `cipher test` in your terminal. It sends test requests through the security filter and test sandbox, and prints a checklist confirming every part of the system responded correctly.
