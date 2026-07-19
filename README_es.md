# 🐚 Shelloma

> **Traductor de Lenguaje Natural a Comandos Terminal para Linux, macOS y Windows impulsado por Ollama local.**

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/moisesfilho/shelloma)
[![Ollama](https://img.shields.io/badge/LLM-Ollama-black?logo=ollama)](https://ollama.com)

**[English](README.md) | [Português](README_pt.md) | [Español](README_es.md)**

**Shelloma** es una aplicación CLI nativa y ultrarrápida desarrollada en Go que traduce tus instrucciones en lenguaje natural directamente a comandos de terminal (Bash/Zsh en Linux/macOS, PowerShell/CMD en Windows) ejecutables en tu equipo, utilizando tu propio modelo local mediante **Ollama**.

---

## 🛠️ ¿Cómo Funciona?

1. **Interpretación**: Shelloma captura tu instrucción en lenguaje natural, detecta automáticamente el contexto de tu sistema (SO, distribución, shell, usuario actual, directorio de trabajo) y envía el contexto a tu API de Ollama local.
2. **Generación de Comando**: El modelo genera estrictamente el comando ejecutable correspondiente (Bash/Zsh para Linux/macOS o PowerShell/CMD para Windows).
3. **Menú Interactivo**: Shelloma muestra el comando formateado en una tarjeta visual y te permite elegir si deseas **Ejecutar**, **Solicitar Explicación del Comando**, **Modificar**, **Copiar al Portapapeles** o **Salir**.
4. **Análisis de Ejecución y Recuperación de Errores**: Tras la ejecución, Shelloma analiza la salida de la terminal. Si el comando falla, identifica la causa y **sugiere automáticamente un comando de corrección**, volviendo recursivamente al comando anterior una vez solucionado el problema.

---

## ✨ Características

- 🚀 **Nativo y Ligero**: Ejecutable único compilado en Go para Linux, macOS y Windows.
- 🔒 **100% Privado y Offline**: Todos tus datos y comandos permanecen localmente en tu equipo.
- 🌐 **Internacionalización Nativa (i18n)**: Soporte completo para **Inglés (`en`)**, **Portugués Brasileño (`pt`)** y **Español (`es`)** mediante archivos JSON incrustados.
- 🤖 **Selección Automática de Modelo**: Detecta los modelos instalados en Ollama y selecciona automáticamente el mejor modelo orientado a código y terminal disponible.
- 💡 **Recuperación de Ollama Offline**: Detecta cuando Ollama está fuera de línea y te ofrece de forma interactiva iniciar el servicio según tu SO (`ollama serve`, `brew services start ollama`, `sudo systemctl start ollama`).
- 📦 **Distribución Multiplataforma**: Ejecutables para Windows (`.exe`), macOS y Linux (`.deb`, `.rpm`, AppImage, Flatpak).

---

## ⚙️ Dependencias

1. **Sistema Operativo**: Linux, macOS o Windows.
2. **Ollama**: Servicio de Ollama instalado y activo localmente.
3. **Go 1.20+** *(opcional)*: Solo necesario si deseas compilar la aplicación desde el código fuente.

---

## 🦙 Cómo Instalar y Configurar Ollama

**[Ollama](https://ollama.com)** es la herramienta de código abierto utilizada para ejecutar modelos de inteligencia artificial localmente en tu ordenador.

- **Sitio Web Oficial**: [https://ollama.com](https://ollama.com)
- **Repositorio GitHub**: [https://github.com/ollama/ollama](https://github.com/ollama/ollama)

### Instalación en Linux

Ejecuta el comando oficial de instalación en tu terminal Linux:

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

*Verifica el estado del servicio con `sudo systemctl status ollama` o inícialo con `sudo systemctl start ollama` / `ollama serve`.*

### Instalación en macOS

- **Opción 1 (Homebrew)**:
  ```bash
  brew install ollama
  brew services start ollama
  ```
- **Opción 2 (Instalador Oficial)**: Descarga el archivo `.zip` desde [ollama.com/download/macOS](https://ollama.com/download/Ollama-darwin.zip), descomprime y muévelo a `/Applications`.

### Instalación en Windows

- **Opción 1 (Winget)**:
  ```powershell
  winget install Ollama.Ollama
  ```
- **Opción 2 (Instalador Oficial)**: Descarga el instalador `OllamaSetup.exe` desde [ollama.com/download/windows](https://ollama.com/download/OllamaSetup.exe) y ejecuta la instalación.

*Tras la instalación en Windows, Ollama permanecerá activo en la bandeja del sistema (System Tray) o podrá iniciarse desde la terminal mediante `ollama serve`.*

---

## 🎯 Modelos Recomendados para Shell y Código

Para obtener la mejor generación de comandos de terminal, recomendamos modelos optimizados para código. Ejecuta `ollama pull <modelo>` para descargar el modelo deseado:

1. **Qwen 2.5 Coder 1.5B** *(Altamente Recomendado - Ligero y Rápido)*:
   ```bash
   ollama pull qwen2.5-coder:1.5b
   ```
2. **Qwen 2.5 Coder 7B** *(Gran precisión para tareas complejas)*:
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

## 📦 Guía de Instalación y Descargas

Los binarios y paquetes precompilados de Shelloma están disponibles para su descarga en la página oficial de **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** (`.deb`, `.rpm`, `.tar.gz`, `.zip`, `AppImage`, `.flatpak`, `.exe`).

### Opción 1: Descarga de Paquetes Precompilados (Recomendado)

Visita **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** y elige el formato adecuado para tu sistema operativo:

#### 🐧 Linux
- **Debian / Ubuntu / Linux Mint / Pop!_OS (`.deb`)**:
  ```bash
  sudo dpkg -i shelloma_1.1.0_amd64.deb
  ```
  *(Durante la instalación del paquete `.deb`, se mostrará un asistente interactivo para elegir el idioma).*

- **Fedora / RedHat / CentOS / RHEL (`.rpm`)**:
  ```bash
  sudo rpm -i shelloma_1.1.0_amd64.rpm
  ```

- **Flatpak (Universal)**:
  ```bash
  flatpak install Shelloma-x86_64.flatpak
  ```

- **AppImage (Universal - Binario portátil sin necesidad de instalación)**:
  ```bash
  chmod +x Shelloma-x86_64.AppImage
  ./Shelloma-x86_64.AppImage "listar archivos del directorio"
  ```

- **Archivo Comprimido Tarball (`.tar.gz`)**:
  Descarga el `.tar.gz` (`amd64` o `arm64`), extrae el binario y muévelo a tu `PATH` (ej: `~/.local/bin/`).

#### 🍏 macOS
- **Binario Nativo (Intel y Apple Silicon M1/M2/M3)**:
  1. Descarga `shelloma_1.1.0_darwin_arm64.tar.gz` (Apple Silicon) o `shelloma_1.1.0_darwin_amd64.tar.gz` (Intel) en Releases.
  2. Extrae el paquete y hazlo ejecutable:
     ```bash
     tar -xzf shelloma_1.1.0_darwin_arm64.tar.gz
     chmod +x shelloma
     sudo mv shelloma /usr/local/bin/
     ```

#### 🪟 Windows
- **Ejecutable Nativo (`.exe` / `.zip`)**:
  1. Descarga `shelloma_1.1.0_windows_amd64.zip` (64-bit) o `shelloma_1.1.0_windows_arm64.zip` (ARM64) desde Releases.
  2. Extrae el archivo zip.
  3. Mueve `shelloma.exe` a una carpeta de tu sistema (ej: `C:\Program Files\Shelloma` o `C:\Tools\`).
  4. *(Opcional)* Añade la ruta del directorio a las Variables de Entorno del sistema (`PATH`) para ejecutar `shelloma` directamente desde cualquier PowerShell o Símbolo del sistema (CMD).

### Opción 2: Compilar desde el Código Fuente

```bash
# 1. Clonar el repositorio
git clone https://github.com/moisesfilho/shelloma.git
cd shelloma

# 2. Compilar e instalar en el directorio del usuario (~/.local/bin)
make build
make install-user
```

---

## 🚀 Guía de Uso Simplificada

### Uso Básico

Ejecuta `shelloma` seguido de tu instrucción entre comillas:

```bash
# Ejemplos en Español
shelloma "listar todos los archivos pdf en la carpeta descargas"
shelloma "verificar el espacio en disco y uso de memoria"
shelloma "crear una carpeta fotos y mover todas las imagenes png a ella"

# Ejemplos con banderas de idioma
shelloma -l en "list all active docker containers"
shelloma -l pt "listar todos os arquivos da pasta atual"
```

### Menú Interactivo en CLI

Al generar un comando, Shelloma mostrará la tarjeta y esperará tu elección:

```text
┌────────────────────────────────────────────┐
│  ls -la ~/Downloads/*.pdf                 │
└────────────────────────────────────────────┘

Opciones: [Enter/y: Ejecutar] [e: Explicar] [m: Modificar] [c: Copiar] [q/n: Salir]:
```

- **Presionar `Enter` o `y`**: Ejecuta el comando directamente en la terminal.
- **Escribir `e`**: Solicita a Ollama una explicación detallada línea por línea.
- **Escribir `m`**: Abre un aviso para editar el comando antes de ejecutarlo.
- **Escribir `c`**: Copia el comando directamente al portapapeles del sistema.
- **Escribir `q` o `n`**: Cancela la operación de forma segura.

### Configuración y Utilidades

```bash
# Mostrar la configuración actual
shelloma config

# Cambiar el idioma guardado por defecto (en, pt, es)
shelloma config set lang es

# Definir un modelo específico de Ollama
shelloma config set model qwen2.5-coder:1.5b

# Listar los modelos de Ollama instalados
shelloma models
```

### Banderas de Comando

```text
Opciones:
  -l, --lang string    Idioma: en, pt, es (predeterminado: en)
  -m, --model string   Modelo Ollama a utilizar (ej: qwen2.5-coder:1.5b)
  -url string          URL de la API de Ollama (predeterminado: http://localhost:11434)
  -y, --yes            Ejecutar el comando generado automáticamente sin meú
  -v, --version        Mostrar la versión de Shelloma
```

---

## 🧪 Desarrollo y Calidad de Código

Shelloma adopta prácticas estrictas de **Clean Code**, el **Principio de Responsabilidad Única (SRP)** y una arquitectura organizada por funcionalidad (**Package by Feature**):

- **`pkg/cli`**: Orquestación de línea de comandos, lectura de flags y flujo interactivo.
- **`pkg/ui`**: Interfaz visual de terminal, renderizado de cards, estilos ANSI y portapapeles.
- **`pkg/ollama`**: Cliente de la API de Ollama, limpieza de prompts y diagnósticos de ejecución.
- **`pkg/sysinfo`**: Detección de sistema operativo, distribución Linux, shell y usuario.
- **`pkg/config`**: Gestión y persistencia de las configuraciones del usuario.
- **`pkg/i18n`**: Soporte nativo e integrado para internacionalización (inglés, portugués y español).

### Análisis Estático y Automatización

```bash
# Ejecutar suite de pruebas unitarias
make test

# Ejecutar análisis estático de código (golangci-lint / staticcheck / go vet)
make lint

# Compilar binario local (ejecuta lint y test automáticamente antes de compilar)
make build
```

El proyecto incluye integración con **`golangci-lint`** y un **Git Pre-Commit Hook** (`.git/hooks/pre-commit`) que ejecuta automáticamente la verificación estática de código y las pruebas unitarias en cada commit y compilación.

---

## 📄 Licencia

Este proyecto está licenciado bajo la licencia **MIT** - consulte el archivo [LICENSE](LICENSE) para obtener más detalles.
