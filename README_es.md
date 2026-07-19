# 🐚 Shelloma

> **Traductor de Lenguaje Natural a Comandos Shell para Linux (Debian/Ubuntu) impulsado por Ollama local.**

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20(Debian%2FUbuntu)-E6007E?logo=debian)](https://www.debian.org)
[![Ollama](https://img.shields.io/badge/LLM-Ollama-black?logo=ollama)](https://ollama.com)

**[English](README.md) | [Português](README_pt.md) | [Español](README_es.md)**

**Shelloma** es una aplicación CLI nativa y ultrarrápida desarrollada en Go que traduce tus instrucciones en lenguaje natural directamente a comandos de terminal (Bash/Zsh) ejecutables en Linux, utilizando tu propio modelo local mediante **Ollama**.

---

## 🛠️ ¿Cómo Funciona?

1. **Interpretación**: Shelloma captura tu instrucción en lenguaje natural, detecta automáticamente el contexto de tu sistema Linux (distribución, versión, usuario actual, directorio de trabajo) y envía el contexto a tu API de Ollama local.
2. **Generación de Comando**: El modelo genera estrictamente el comando ejecutable Bash/Zsh correspondiente.
3. **Menú Interactivo**: Shelloma muestra el comando formateado en una tarjeta visual y te permite elegir si deseas **Ejecutar**, **Solicitar Explicación del Comando**, **Modificar**, **Copiar al Portapapeles** o **Salir**.
4. **Análisis de Ejecución y Recuperación de Errores**: Tras la ejecución, Shelloma analiza la salida de la terminal. Si el comando falla, identifica la causa y **sugiere automáticamente un comando de corrección**, volviendo recursivamente al comando anterior una vez solucionado el problema.

---

## ✨ Características

- 🚀 **Nativo y Ligero**: Ejecutable único compilado en Go, sin necesidad de Python, Node.js ni runtimes externos.
- 🔒 **100% Privado y Offline**: Todos tus datos y comandos permanecen localmente en tu equipo.
- 🌐 **Internacionalización Nativa (i18n)**: Soporte completo para **Inglés (`en`)**, **Portugués Brasileño (`pt`)** y **Español (`es`)** mediante archivos JSON incrustados.
- 🤖 **Selección Automática de Modelo**: Detecta los modelos instalados en Ollama y selecciona automáticamente el mejor modelo orientado a código y terminal disponible.
- 💡 **Recuperación de Ollama Offline**: Detecta cuando Ollama está fuera de línea y te ofrece de forma interactiva iniciar el servicio (`sudo systemctl start ollama`) con un bucle de reintento automático.
- 📦 **Instalación Simplificada en `.deb`**: Instalador Debian listo con un asistente interactivo de selección de idioma durante `sudo dpkg -i`.

---

## ⚙️ Dependencias

1. **Sistema Operativo**: Linux (Debian, Ubuntu o distribuciones derivadas).
2. **Ollama**: Servicio de Ollama instalado y activo localmente.
3. **Go 1.20+** *(opcional)*: Solo necesario si deseas compilar la aplicación desde el código fuente.

---

## 🦙 Cómo Instalar y Configurar Ollama

**[Ollama](https://ollama.com)** es la herramienta de código abierto utilizada para ejecutar modelos de inteligencia artificial localmente en tu ordenador.

- **Sitio Web Oficial**: [https://ollama.com](https://ollama.com)
- **Repositorio GitHub**: [https://github.com/ollama/ollama](https://github.com/ollama/ollama)

### Paso 1: Instalar Ollama en Linux

Ejecuta el comando oficial de instalación en tu terminal Linux:

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

### Paso 2: Verificar el Estado del Servicio

Generalmente Ollama se inicia automáticamente como demonio del sistema. Puedes verificar su estado con:

```bash
sudo systemctl status ollama
```

*Si el servicio no está en ejecución, inícialo con `sudo systemctl start ollama` o ejecuta manualmente `ollama serve`.*

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

Los binarios y paquetes precompilados de Shelloma están disponibles para su descarga en la página oficial de **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** (`.deb`, `.rpm`, `.tar.gz`, `.snap`, `AppImage`, `.flatpak`).

### Opción 1: Descarga de Paquetes Precompilados (Recomendado)

Visita **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** y elige el formato adecuado para tu distribución Linux:

- **Debian / Ubuntu / Linux Mint / Pop!_OS (`.deb`)**:
  ```bash
  sudo dpkg -i shelloma_1.0.1_amd64.deb
  ```
  *(Durante la instalación del paquete `.deb`, se mostrará un asistente interactivo para elegir el idioma).*

- **Fedora / RedHat / CentOS / RHEL (`.rpm`)**:
  ```bash
  sudo rpm -i shelloma_1.0.1_amd64.rpm
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

## 📄 Licencia

Este proyecto está bajo la Licencia **MIT** - consulta el archivo [LICENSE](LICENSE) para más detalles.
