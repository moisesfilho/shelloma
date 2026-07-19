# 🐚 Shelloma

> **Natural Language to Terminal Command Translator for Linux, macOS & Windows powered by local Ollama.**

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/moisesfilho/shelloma)
[![Ollama](https://img.shields.io/badge/LLM-Ollama-black?logo=ollama)](https://ollama.com)

**[English](README.md) | [Português](README_pt.md) | [Español](README_es.md)**

**Shelloma** is a native, ultra-fast CLI application written in Go that translates natural language instructions directly into executable terminal commands (Bash/Zsh on Linux/macOS, PowerShell/CMD on Windows) using your local LLM models via **Ollama**.

---

## 🛠️ How It Works

1. **Interpretation**: Shelloma captures your natural language instruction, automatically detects your system context (OS, distribution, shell, current user, working directory), and sends the context to your local Ollama API.
2. **Command Generation**: The model generates strictly the corresponding executable terminal command (Bash/Zsh for Linux/macOS or PowerShell/CMD for Windows).
3. **Interactive Menu**: Shelloma presents the generated command in a formatted visual card and lets you choose whether to **Execute**, **Request Command Explanation**, **Modify**, **Copy to Clipboard**, or **Quit**.
4. **Execution Analysis & Error Recovery**: After execution, Shelloma analyzes the terminal output. If the command fails, it identifies the root cause and **automatically suggests a fix command**, recursively returning to the previous command once the issue is resolved.

---

## ✨ Features

- 🚀 **Native & Lightweight**: Single self-contained binary compiled in Go for Linux, macOS, and Windows.
- 🔒 **100% Private & Offline**: All data and commands stay on your local machine.
- 🌐 **Native Internationalization (i18n)**: Full support for **English (`en`)**, **Brazilian Portuguese (`pt`)**, and **Spanish (`es`)** via embedded JSON files (`embed.FS`).
- 🤖 **Automatic Model Selection**: Automatically detects installed Ollama models and selects the best available coding/shell model.
- 💡 **Ollama Offline Recovery**: Detects when Ollama is offline and interactively offers to start the service based on your OS (`ollama serve`, `brew services start ollama`, `sudo systemctl start ollama`).
- 📦 **Cross-Platform Distribution**: Binaries for Windows (`.exe`), macOS, and Linux (`.deb`, `.rpm`, AppImage, Flatpak).

---

## ⚙️ Dependencies

1. **Operating System**: Linux, macOS, or Windows.
2. **Ollama**: Ollama service installed and running locally.
3. **Go 1.20+** *(optional)*: Only required if compiling from source code.

---

## 🦙 How to Install and Configure Ollama

**[Ollama](https://ollama.com)** is the open-source engine used to run AI models locally on your machine.

- **Official Website**: [https://ollama.com](https://ollama.com)
- **GitHub Repository**: [https://github.com/ollama/ollama](https://github.com/ollama/ollama)

### Installation on Linux

Run the official installation script in your Linux terminal:

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

*Check service status with `sudo systemctl status ollama` or start it with `sudo systemctl start ollama` / `ollama serve`.*

### Installation on macOS

- **Option 1 (Homebrew)**:
  ```bash
  brew install ollama
  brew services start ollama
  ```
- **Option 2 (Official Installer)**: Download the `.zip` from [ollama.com/download/macOS](https://ollama.com/download/Ollama-darwin.zip), unzip, and move it to `/Applications`.

### Installation on Windows

- **Option 1 (Winget)**:
  ```powershell
  winget install Ollama.Ollama
  ```
- **Option 2 (Official Installer)**: Download `OllamaSetup.exe` from [ollama.com/download/windows](https://ollama.com/download/OllamaSetup.exe) and run the installer.

*After installing on Windows, Ollama will run in the System Tray or can be started from the terminal using `ollama serve`.*

---

## 🎯 Recommended Models for Shell & Coding

For optimal terminal command generation, we recommend code-focused models. Run `ollama pull <model>` to download a model:

1. **Qwen 2.5 Coder 1.5B** *(Highly Recommended - Lightweight & Fast)*:
   ```bash
   ollama pull qwen2.5-coder:1.5b
   ```
2. **Qwen 2.5 Coder 7B** *(High accuracy for complex tasks)*:
   ```bash
   ollama pull qwen2.5-coder:7b
   ```
3. **DeepSeek Coder 6.7B**:
   ```bash
   ollama pull deepseek-coder:6.7b
   ```
4. **Llama 3.2 3B**:
   ```bash
   ollama pull llama3.2:3b
   ```

---

## 📦 Installation Guide & Downloads

Pre-compiled binaries and packages for Shelloma are available for download on the official **[GitHub Releases Page](https://github.com/moisesfilho/shelloma/releases)** (`.deb`, `.rpm`, `.tar.gz`, `.zip`, `AppImage`, `.flatpak`, `.exe`).

### Option 1: Pre-compiled Packages & Downloads (Recommended)

Visit **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** and choose the right format for your operating system:

#### 🐧 Linux
- **Debian / Ubuntu / Linux Mint / Pop!_OS (`.deb`)**:
  ```bash
  sudo dpkg -i shelloma_1.0.2_amd64.deb
  ```
  *(During `.deb` installation, an interactive prompt asks for your preferred default language).*

- **Fedora / RedHat / CentOS / RHEL (`.rpm`)**:
  ```bash
  sudo rpm -i shelloma_1.0.2_amd64.rpm
  ```

- **Flatpak (Universal)**:
  ```bash
  flatpak install Shelloma-x86_64.flatpak
  ```

- **AppImage (Universal - Standalone portable binary)**:
  ```bash
  chmod +x Shelloma-x86_64.AppImage
  ./Shelloma-x86_64.AppImage "list files in directory"
  ```

- **Tarball Archive (`.tar.gz`)**:
  Download `.tar.gz` (`amd64` or `arm64`), extract the binary, and place it in your `PATH` (e.g., `~/.local/bin/`).

#### 🍏 macOS
- **Native Binary (Intel & Apple Silicon M1/M2/M3)**:
  1. Download `shelloma_1.0.2_darwin_arm64.tar.gz` (Apple Silicon) or `shelloma_1.0.2_darwin_amd64.tar.gz` (Intel) from Releases.
  2. Extract and make it executable:
     ```bash
     tar -xzf shelloma_1.0.2_darwin_arm64.tar.gz
     chmod +x shelloma
     sudo mv shelloma /usr/local/bin/
     ```

#### 🪟 Windows
- **Native Executable (`.exe` / `.zip`)**:
  1. Download `shelloma_1.0.2_windows_amd64.zip` (64-bit) or `shelloma_1.0.2_windows_arm64.zip` (ARM64) from Releases.
  2. Extract the zip archive.
  3. Move `shelloma.exe` to a folder on your system (e.g., `C:\Program Files\Shelloma` or `C:\Tools\`).
  4. *(Optional)* Add the directory path to system Environment Variables (`PATH`) to run `shelloma` directly from any PowerShell or Command Prompt (CMD).

### Option 2: Build from Source

```bash
# 1. Clone repository
git clone https://github.com/moisesfilho/shelloma.git
cd shelloma

# 2. Build and install to user binary directory (~/.local/bin)
make build
make install-user
```

---

## 🚀 Simplified Usage Guide

### Basic Usage

Run `shelloma` followed by your instruction in quotes:

```bash
# English examples
shelloma "list all pdf files in downloads folder"
shelloma "check disk space and memory usage"
shelloma "create a folder named images and move all png files into it"

# Multi-language flags
shelloma -l pt "listar todos os arquivos da pasta atual"
shelloma -l es "mostrar procesos que consumen mas cpu"
```

### Interactive CLI Menu

When a command is generated, Shelloma displays the card and waits for your choice:

```text
┌────────────────────────────────────────────┐
│  ls -la ~/Downloads/*.pdf                 │
└────────────────────────────────────────────┘

Options: [Enter/y: Execute] [e: Explain] [m: Modify] [c: Copy] [q/n: Quit]:
```

- **Press `Enter` or `y`**: Executes the command directly in terminal.
- **Type `e`**: Requests a line-by-line explanation from Ollama.
- **Type `m`**: Opens a prompt to edit the command before execution.
- **Type `c`**: Copies the command directly to your system clipboard.
- **Type `q` or `n`**: Cancels the operation safely.

### Configuration & Utilities

```bash
# Display current configuration
shelloma config

# Change default saved language (en, pt, es)
shelloma config set lang en

# Set a specific Ollama model
shelloma config set model qwen2.5-coder:1.5b

# List installed Ollama models
shelloma models
```

### Command Flags

```text
Options:
  -l, --lang string    Language: en, pt, es (default: en)
  -m, --model string   Ollama model to use (e.g. qwen2.5-coder:1.5b)
  -url string          Ollama API URL (default: http://localhost:11434)
  -y, --yes            Execute generated command automatically without confirmation
  -v, --version        Display Shelloma version
```

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.
