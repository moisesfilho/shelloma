# 🐚 Shelloma

> **Natural Language to Shell Translator for Linux (Debian/Ubuntu) powered by Ollama.**

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20(Debian%2FUbuntu)-E6007E?logo=debian)](https://www.debian.org)

**Shelloma** is a native, ultra-fast Linux CLI application written in Go that translates natural language instructions into executable Bash/Zsh terminal commands using your local Ollama LLM.

---

## ✨ Features

- 🚀 **Native & Lightweight**: Single self-contained binary written in Go with no Python or Node.js runtime required.
- 🔒 **100% Local & Private**: Runs completely offline using local models via Ollama.
- 🌐 **Multi-language Support (i18n)**: Native translation in **English (`en`)**, **Português Brasileiro (`pt`)**, and **Español (`es`)** loaded via embedded JSON files (`embed.FS`).
- 🤖 **Automatic Model Selection**: Automatically detects and picks the best installed coding model (e.g., `qwen2.5-coder`, `deepseek-coder`, `llama3.2`).
- 🛡️ **Execution Analysis & Auto-Recovery**: Analyzes terminal output after command execution. If a command fails, Shelloma analyzes the error and suggests a fix, recursively returning to the previous command upon success.
- 💡 **Offline Ollama Recovery**: Detects when Ollama is offline and interactively offers to start the service (`sudo systemctl start ollama`) with an automatic retry loop.
- 📦 **Debian Package Integration**: Pre-packaged `.deb` file with interactive language setup during `sudo dpkg -i`.

---

## 📦 Installation

### Option 1: Debian / Ubuntu Package (`.deb`)

Download or build the latest `.deb` package and install:

```bash
sudo dpkg -i shelloma_1.0.0_amd64.deb
```

During installation, you will be prompted to select your default language (`en`, `pt`, or `es`).

### Option 2: Build from Source

Prerequisites: Go 1.20+ and `make`.

```bash
git clone https://github.com/moisesfilho/shelloma.git
cd shelloma
make build
make install-user
```

---

## 🚀 Usage

### Basic Usage

Simply pass your instruction in natural language:

```bash
# English
shelloma "list all pdf files in downloads folder"

# Português
shelloma "listar todos os arquivos pdf na pasta downloads"

# Español
shelloma "mostrar uso de memoria y disco"
```

### Interactive Menu

When a command is generated, Shelloma presents a clean, formatted card with options:

```text
┌────────────────────────────────────────────┐
│  ls -la ~/Downloads/*.pdf                 │
└────────────────────────────────────────────┘

Options: [Enter/y: Execute] [e: Explain] [m: Modify] [c: Copy] [q/n: Quit]:
```

- `Enter` or `y`: Execute the command.
- `e`: Request a detailed explanation of the command from Ollama.
- `m`: Edit the command interactively in terminal before running.
- `c`: Copy the command directly to system clipboard.
- `q` or `n`: Cancel operation.

---

## ⚙️ Configuration & Commands

### Change Configuration

```bash
# Show current configuration
shelloma config

# Set default language (en, pt, es)
shelloma config set lang pt

# Set specific Ollama model
shelloma config set model qwen2.5-coder:1.5b

# List installed Ollama models
shelloma models
```

### Command Flags

```text
Usage:
  shelloma "instruction"
  shelloma [options] "instruction"

Options:
  -l, --lang string    Language: en, pt, es (default: en)
  -m, --model string   Ollama model to use (e.g. qwen2.5-coder:1.5b)
  -url string          Ollama API URL (default: http://localhost:11434)
  -y, --yes            Execute generated command automatically without confirmation
  -v, --version        Display Shelloma version
```

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
