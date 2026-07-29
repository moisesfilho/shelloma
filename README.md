# 🐚 Shelloma

> **Natural Language to Terminal Command Translator for Linux, macOS & Windows powered by local Ollama.**

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![CI](https://github.com/moisesfilho/shelloma/actions/workflows/ci.yml/badge.svg)](https://github.com/moisesfilho/shelloma/actions/workflows/ci.yml)
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
- 🌐 **Native Internationalization (i18n)**: Full support for **English (`en`)**, **Brazilian Portuguese (`pt`)**, and **Spanish (`es`)** via embedded JSON files (with translated header and messages).
- 🤖 **Automatic Model Selection**: Automatically detects installed Ollama models and selects the best available coding/shell model.
- 💡 **Ollama Offline Recovery**: Detects when Ollama is offline and interactively offers to start the service based on your OS.
- 🖥️ **Adaptive Interface with Fixed Footer**: Only the footer legend remains fixed at the bottom of the terminal window, without interrupting the natural scrolling of the header and active prompt workspace.
- 🔄 **Interactive Command Refinement**: Allows adjusting the suggested command by selecting `[r: Refine]`, opening an interactive prompt to supply additional feedback or changes.
- 🧭 **Interactive History with Persistence**: All interactive input prompts use a raw terminal reader that enables history navigation with `Up`/`Down` arrows and cursor movement with `Left`/`Right` arrows. History is persisted across sessions (up to 1000 entries).
- 🖥️ **Desktop Terminal Mode**: Run with `--desktop` flag to keep the terminal open after execution ("Press Enter to exit..."). Ideal for launchers and `.desktop` file integration.
- ♾️ **Continuous Command Loop**: After a successful execution, Shelloma prompts for a new instruction, allowing a seamless multi-command workflow without restarting.
- ⌨️ **Rich Keyboard Shortcuts**: The interactive menu supports multiple aliases per action (e.g. `y`/`sim`/`yes` for execute, `e`/`ex`/`explain` for explain), plus `Esc` quits from any prompt.
- 📦 **Cross-Platform Distribution**: Binaries for Windows (`.exe`), macOS, and Linux (`.deb`, `.rpm`, Snap, AppImage, Flatpak). Includes native Linux desktop integration (generating a `.desktop` launcher and registering the official SVG icon).
- ⛓️ **Multi-Step Command Execution**: Executes multi-step command sequences step-by-step with confirmation prompts and logs each step.
- 📝 **Bash Completion**: Automatic shell completion scripts for Debian/Ubuntu packages.

---

## ⚙️ Dependencies

1. **Operating System**: Linux, macOS, or Windows.
2. **Ollama**: Ollama service installed and running locally.
3. **Go 1.23+** *(optional)*: Only required if compiling from source code.

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

Pre-compiled binaries and packages for Shelloma are available for download on the official **[GitHub Releases Page](https://github.com/moisesfilho/shelloma/releases)** (`.deb`, `.rpm`, `.tar.gz`, `.zip`, `.AppImage`, `.flatpak`, Snap, `.exe`).

### Option 1: Pre-compiled Packages & Downloads (Recommended)

Visit **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** and choose the right format for your operating system:

#### 🐧 Linux
- **Debian / Ubuntu / Linux Mint / Pop!_OS (`.deb`)**:
  ```bash
  sudo dpkg -i shelloma_1.4.0_amd64.deb
  ```
  *(During `.deb` installation, an interactive prompt asks for your preferred default language).*

- **Fedora / RedHat / CentOS / RHEL (`.rpm`)**:
  ```bash
  sudo rpm -i shelloma_1.4.0_amd64.rpm
  ```

- **Flatpak (Universal)**:
  ```bash
  flatpak install Shelloma-x86_64.flatpak
  ```

- **Snap (Universal)**:
  ```bash
  snap install shelloma
  ```
  *(Or download `Shelloma-x86_64.snap` from Releases and install with `snap install --dangerous Shelloma-x86_64.snap`)*

- **AppImage (Universal - Standalone portable binary)**:
  ```bash
  chmod +x Shelloma-x86_64.AppImage
  ./Shelloma-x86_64.AppImage "list files in directory"
  ```

- **Tarball Archive (`.tar.gz`)**:
  Download `.tar.gz` (`amd64` or `arm64`), extract the binary, and place it in your `PATH` (e.g., `~/.local/bin/`).

#### 🍏 macOS
- **Native Binary (Intel & Apple Silicon M1/M2/M3)**:
  1. Download `shelloma_1.4.0_darwin_arm64.tar.gz` (Apple Silicon) or `shelloma_1.4.0_darwin_amd64.tar.gz` (Intel) from Releases.
  2. Extract and make it executable:
     ```bash
     tar -xzf shelloma_1.4.0_darwin_arm64.tar.gz
     chmod +x shelloma
     sudo mv shelloma /usr/local/bin/
     ```

#### 🪟 Windows
- **Native Executable (`.exe` / `.zip`)**:
  1. Download `shelloma_1.4.0_windows_amd64.zip` (64-bit) or `shelloma_1.4.0_windows_arm64.zip` (ARM64) from Releases.
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
Options: [Enter/y: Execute] [e: Explain] [m: Modify] [c: Copy] [p: New Prompt] [r: Refine] [a: Adjust Prompt] [q/n: Quit]:
```

- **Press `Enter` or `y`**: Executes the command directly in the terminal.
- **Type `e`**: Requests a detailed, line-by-line explanation of the command from Ollama.
- **Type `m`**: Opens a prompt to manually edit the command before execution.
- **Type `c`**: Copies the command directly to your system clipboard.
- **Type `p`**: Prompts for a completely new command instruction without quitting Shelloma.
- **Type `r`**: Refines the suggested command by opening an interactive prompt to supply additional feedback or adjustments for the LLM to generate a revised suggestion.
- **Type `a`**: Adjusts/edits the original natural language query prompt that triggered the suggestion.
- **Type `q` or `n`**: Cancels the operation and exits safely.
- **Press `Esc`**: Quits immediately from any prompt or the interactive menu.

### Configuration & Utilities

```bash
# Display current configuration
shelloma config

# Display the compact configuration schema
shelloma config docs

# Change default saved language (en, pt, es)
shelloma config set lang en

# Set a specific Ollama model
shelloma config set model qwen2.5-coder:1.5b

# Set model sampling temperature (default: 0.1)
shelloma config set temperature 0.5

# Enable/Disable automatic command execution without confirmation
shelloma config set auto_execute true

# List installed Ollama models
shelloma models

# Open application execution logs
shelloma logs
```

### 🤖 Natural Language Configuration

Shelloma allows you to request configuration updates using natural language! The local AI parses the request and generates the matching Shelloma CLI command to be executed:

- `shelloma "change shelloma model to llama3"` ➔ Suggests: `shelloma config set model llama3`
- `shelloma "set language to spanish"` ➔ Suggests: `shelloma config set lang es`
- `shelloma "set temperature to 0.8"` ➔ Suggests: `shelloma config set temperature 0.8`
- `shelloma "enable auto execution"` ➔ Suggests: `shelloma config set auto_execute true`

### 📝 Custom Translation Rules

You can define custom formatting and execution preferences that will be injected directly into the Ollama system prompt to guide command generation (e.g. setting default folders, choosing specific CLI text editors or programs to open specific file extensions, etc.).

Commands:
```bash
# Add a new custom rule
shelloma rules add "Always open images with xdg-open"

# List all saved rules
shelloma rules list

# Edit a saved rule by index
shelloma rules edit 1 "Always open images with feh"

# Delete a saved rule by index
shelloma rules delete 1
```

### 📋 Structured Logging

Every command suggestion and execution outcome is automatically recorded in a structured log file (using the standard cache directory for each operating system: `~/.cache/shelloma/shelloma.log` on Linux/macOS, or `%LocalAppData%\shelloma\shelloma.log` on Windows).

Run `shelloma logs` to open it. You can choose to:
1. Print the logs in a clean structured format directly in your terminal.
2. Open the raw log file in your system's default text editor.

### ⛓️ Multi-Step Command Execution

When a user query requires multiple execution steps, Shelloma instructs the local model to structure the commands as separate command lines (one per line).

Features:
- **Confirmation Prompts**: Shelloma presents the planned command sequence and asks for user confirmation for each step individually.
- **Directory Traversal Persistence**: `cd` commands are intercepted and applied to the process workspace, ensuring subsequent execution steps run in the correct context.
- **Error Analysis & Logs**: Each step is checked for errors and recorded as an individual execution log entry.

### ⚠️ Dangerous Commands Protection

To prevent accidental system damage, Shelloma includes a security validation for potentially dangerous commands across Linux, macOS, and Windows (e.g., `rm`, `dd`, `mkfs`, `shred`, `chmod`, `chown`, `Remove-Item`, `del`, `rd`, `rmdir`, `format`, `Format-Volume`).

- **Alerts**: When a dangerous command is suggested, a visual warning alert is displayed immediately below the command card.
- **Safety Word Confirmation**: If you attempt to execute a dangerous command, you will be prompted to type the security word `"CONFIRM"` (case-sensitive) to proceed. If typed incorrectly, execution is aborted.
- **Configurable List**: The list of dangerous commands is fully editable.
- **Disable Safety Checks**: You can bypass this validation entirely if desired.

#### Configuration Commands:

```bash
# Add/change the list of dangerous commands (comma-separated list)
shelloma config set dangerous "rm,dd,mkfs,shred,chmod,chown,Remove-Item,del,rd,rmdir,format,Format-Volume"

# Disable the dangerous command protection check entirely
shelloma config set disable_dangerous_check true

# Enable the check again (default)
shelloma config set disable_dangerous_check false
```

### Command Flags

```text
Options:
  -l, --lang string    Language: en, pt, es (default: en)
  -m, --model string   Ollama model to use (e.g. qwen2.5-coder:1.5b)
  -url string          Ollama API URL (default: http://localhost:11434)
  -y, --yes            Execute generated command automatically without confirmation
      --desktop        Desktop terminal mode (pauses with "Press Enter to exit...")
  -v, --version        Display Shelloma version
```

> **Desktop Mode**: When launched via a `.desktop` file or GUI launcher, use `--desktop` so the terminal stays open after execution — useful for seeing output before the window closes.

---

## 🧪 Development & Code Quality

Shelloma follows strict **Clean Code** principles, the **Single Responsibility Principle (SRP)**, and a **Package by Feature** architecture:

- **`pkg/cli`**: Command-line orchestration, flag parsing, and interactive loop.
- **`pkg/ui`**: Terminal visual UI, card rendering, ANSI styling, and clipboard integration.
- **`pkg/ollama`**: Ollama API client, prompt cleaning, and execution diagnostics.
- **`pkg/sysinfo`**: Operating system, Linux distribution, shell, and user detection.
- **`pkg/config`**: User configuration management and persistence.
- **`pkg/i18n`**: Native embedded internationalization (English, Portuguese, Spanish).

### Static Analysis & Automation

```bash
# Run unit test suite
make test

# Run static code analysis (golangci-lint / staticcheck / go vet)
make lint

# Generate application icons (requires Python + Pillow)
make icons

# Build local binary (runs lint, test, icons automatically before compiling)
make build

# Build for all platforms (Linux, macOS, Windows — amd64 + arm64)
make build-all

# Generate .deb, .tar.gz, .AppImage, or .flatpak packages
make deb
make tar
make appimage
make flatpak

# Install binary + .desktop file + icons to ~/.local/bin/
make install-user
```

The project includes integration with **`golangci-lint`**, a **Git Pre-Commit Hook** (`.git/hooks/pre-commit`), and a **CI/CD pipeline** via GitHub Actions:
- **CI Quality Gate** (`.github/workflows/ci.yml`): Runs lint + test on push/PR to `main` and `develop`.
- **Release Pipeline** (`.github/workflows/release.yml`): On tag push, runs lint + test, then builds and publishes all packages (`.deb`, `.rpm`, `.tar.gz`, `.zip`, `.AppImage`, `.flatpak`, Snap) to GitHub Releases.

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.
