# 🐚 Shelloma

> **Tradutor de Linguagem Natural para Comandos Terminal (Linux, macOS e Windows) alimentado por Ollama local.**

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)](https://github.com/moisesfilho/shelloma)
[![Ollama](https://img.shields.io/badge/LLM-Ollama-black?logo=ollama)](https://ollama.com)

**[English](README.md) | [Português](README_pt.md) | [Español](README_es.md)**

O **Shelloma** é uma aplicação CLI nativa e ultrarrápida desenvolvida em Go que traduz suas instruções em linguagem natural diretamente para comandos de terminal (Bash/Zsh no Linux/macOS, PowerShell/CMD no Windows) executáveis no seu computador, utilizando o seu próprio modelo local via **Ollama**.

---

## 🛠️ Como Funciona?

1. **Interpretação**: O Shelloma captura sua instrução em português (ou outro idioma configurado), detecta automaticamente as especificidades do seu sistema operacional (Linux, macOS ou Windows, shell em uso, usuário e diretório atual) e envia o contexto para a API do Ollama local.
2. **Geração de Comando**: O modelo gera estritamente o comando de terminal correspondente (Bash/Zsh para Linux/macOS ou PowerShell/CMD para Windows).
3. **Menu Interativo**: O Shelloma exibe o comando formatado em um card visual e permite que você escolha se deseja **Executar**, **Pedir Explicação do Comando**, **Modificar**, **Copiar para Clipboard** ou **Sair**.
4. **Análise de Execução e Recuperação de Erros**: Após a execução, o Shelloma analisa o retorno do terminal. Se houver falha, ele identifica a causa e **sugere automaticamente um comando de correção**, retornando ao comando anterior de forma encadeada assim que o erro for corrigido.

---

## ✨ Recursos

- 🚀 **Nativo e Leve**: Executável único compilado em Go para Linux, macOS e Windows.
- 🔒 **100% Privado e Offline**: Nenhum dado ou comando sai da sua máquina.
- 🌐 **Internacionalização Nativa (i18n)**: Suporte completo para **Português Brasileiro (`pt`)**, **Inglês (`en`)** e **Espanhol (`es`)** através de arquivos JSON embarcados.
- 🤖 **Seleção Automática de Modelo**: Detecta os modelos instalados no Ollama e seleciona automaticamente o melhor modelo focado em código/shell disponível.
- 💡 **Detecção e Inicialização do Ollama**: Se o serviço Ollama estiver parado, o Shelloma avisa e oferece um comando rápido interativo para iniciá-lo de acordo com seu SO (`ollama serve`, `brew services start ollama`, `sudo systemctl start ollama`).
- 📦 **Instalação Multiplataforma**: Executáveis compilados (`.exe`, binários nativos) e pacotes Linux (`.deb`, `.rpm`, AppImage, Flatpak).

---

## ⚙️ Dependências

1. **Sistema Operacional**: Linux, macOS ou Windows.
2. **Ollama**: Serviço do Ollama instalado e ativo localmente.
3. **Go 1.20+** *(opcional)*: Apenas caso deseje compilar a aplicação a partir do código-fonte.

---

## 🦙 Como Instalar e Configurar o Ollama

O **[Ollama](https://ollama.com)** é a ferramenta aberta responsável por rodar modelos de inteligência artificial localmente no seu computador.

- **Website Oficial**: [https://ollama.com](https://ollama.com)
- **Repositório GitHub**: [https://github.com/ollama/ollama](https://github.com/ollama/ollama)

### Instalação no Linux

No terminal Linux, execute o comando oficial de instalação:

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

*Verifique o status do serviço com `sudo systemctl status ollama` ou inicie com `sudo systemctl start ollama` / `ollama serve`.*

### Instalação no macOS

- **Opção 1 (Homebrew)**:
  ```bash
  brew install ollama
  brew services start ollama
  ```
- **Opção 2 (Instalador Oficial)**: Baixe o arquivo `.zip` em [ollama.com/download/macOS](https://ollama.com/download/Ollama-darwin.zip), descompacte e mova para a pasta `/Applications`.

### Instalação no Windows

- **Opção 1 (Winget)**:
  ```powershell
  winget install Ollama.Ollama
  ```
- **Opção 2 (Instalador Oficial)**: Baixe o instalador `OllamaSetup.exe` em [ollama.com/download/windows](https://ollama.com/download/OllamaSetup.exe) e siga o assistente.

*Após instalar no Windows, o Ollama ficará ativo na bandeja do sistema (System Tray) ou poderá ser iniciado pelo terminal via `ollama serve`.*

---

## 🎯 Modelos Recomendados para Shell e Código

Para obter as melhores respostas de comandos de terminal Linux, recomendamos utilizar modelos treinados para código. Execute o comando `ollama pull <modelo>` para baixar o modelo desejado:

1. **Qwen 2.5 Coder 1.5B** *(Altamente Recomendado - Leve e Ultrarrápido)*:
   ```bash
   ollama pull qwen2.5-coder:1.5b
   ```
2. **Qwen 2.5 Coder 7B** *(Excelente precisão para tarefas complexas)*:
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

## 📦 Guia de Instalação e Downloads

Todos os binários e pacotes pré-compilados do Shelloma estão disponíveis para download na página oficial de **[Releases no GitHub](https://github.com/moisesfilho/shelloma/releases)** (`.deb`, `.rpm`, `.tar.gz`, `.zip`, `AppImage`, `.flatpak`, `.exe`).

### Opção 1: Download de Binários e Pacotes Pré-compilados (Recomendado)

Acesse a página de **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** e escolha a opção ideal para o seu sistema operacional:

#### 🐧 Linux
- **Debian / Ubuntu / Linux Mint / Pop!_OS (`.deb`)**:
  ```bash
  sudo dpkg -i shelloma_1.1.0_amd64.deb
  ```
  *(Durante a instalação do `.deb`, será exibido um assistente no terminal para selecionar o idioma padrão).*

- **Fedora / RedHat / CentOS / RHEL (`.rpm`)**:
  ```bash
  sudo rpm -i shelloma_1.1.0_amd64.rpm
  ```

- **Flatpak (Universal)**:
  ```bash
  flatpak install Shelloma-x86_64.flatpak
  ```

- **AppImage (Universal - Executável portável sem necessidade de instalação)**:
  ```bash
  chmod +x Shelloma-x86_64.AppImage
  ./Shelloma-x86_64.AppImage "listar arquivos do diretorio"
  ```

- **Arquivo Comprimido Tarball (`.tar.gz`)**:
  Baixe o `.tar.gz` (`amd64` ou `arm64`), extraia o binário e mova para seu `PATH` (ex: `~/.local/bin/`).

#### 🍏 macOS
- **Binário Nativo (Intel & Apple Silicon M1/M2/M3)**:
  1. Baixe `shelloma_1.1.0_darwin_arm64.tar.gz` (Apple Silicon) ou `shelloma_1.1.0_darwin_amd64.tar.gz` (Intel) na página de Releases.
  2. Extraia o pacote e torne-o executável:
     ```bash
     tar -xzf shelloma_1.1.0_darwin_arm64.tar.gz
     chmod +x shelloma
     sudo mv shelloma /usr/local/bin/
     ```

#### 🪟 Windows
- **Executável Nativamente Compilado (`.exe` / `.zip`)**:
  1. Baixe `shelloma_1.1.0_windows_amd64.zip` (64-bit) ou `shelloma_1.1.0_windows_arm64.zip` (ARM64) na página de Releases.
  2. Extraia o arquivo zip.
  3. Mova o executável `shelloma.exe` para o diretório desejado (ex: `C:\Program Files\Shelloma` ou `C:\Tools\`).
  4. *(Opcional)* Adicione o diretório às variáveis de ambiente de sistema (`PATH`) para executar `shelloma` diretamente em qualquer terminal PowerShell ou Prompt de Comando (CMD).

### Opção 2: Compilando a partir do Código-Fonte

```bash
# 1. Clonar o repositório
git clone https://github.com/moisesfilho/shelloma.git
cd shelloma

# 2. Compilar e instalar no diretório do usuário (~/.local/bin)
make build
make install-user
```

---

## 🚀 Guia de Uso Simplificado

### Uso Básico

Basta chamar o `shelloma` seguido da sua instrução entre aspas:

```bash
# Exemplos em Português
shelloma "listar todos os arquivos pdf da pasta downloads"
shelloma "verificar espaço em disco e uso de memória"
shelloma "criar uma pasta chamada fotos e mover todas as imagens png para ela"

# Exemplos em outros idiomas
shelloma -l en "list all active docker containers"
shelloma -l es "mostrar procesos que consumen mas cpu"
```

### Menu de Opções da CLI

Ao gerar o comando, o Shelloma exibirá o card e aguardará sua ação:

```text
┌────────────────────────────────────────────┐
│  ls -la ~/Downloads/*.pdf                 │
└────────────────────────────────────────────┘

Opções: [Enter/y: Executar] [e: Explicar] [m: Modificar] [c: Copiar] [q/n: Sair]:
```

- **Pressionar `Enter` ou `y`**: Executa o comando diretamente no terminal.
- **Digite `e`**: Solicita uma explicação linha por linha do comando ao Ollama.
- **Digite `m`**: Abre um prompt para você editar o comando antes de executar.
- **Digite `c`**: Copia o comando diretamente para a área de transferência do sistema.
- **Digite `q` ou `n`**: Cancela a operação de forma segura.

### Comandos de Configuração e Utilitários

```bash
# Exibir as configurações atuais do Shelloma
shelloma config

# Alterar o idioma padrão salvo (pt, en, es)
shelloma config set lang pt

# Definir um modelo específico do Ollama
shelloma config set model qwen2.5-coder:1.5b

# Listar os modelos do Ollama instalados na sua máquina
shelloma models
```

### ⚠️ Proteção de Comandos Perigosos

Para evitar danos acidentais ao sistema, o Shelloma inclui uma validação de segurança para comandos potencialmente perigosos em Linux, macOS e Windows (por exemplo, `rm`, `dd`, `mkfs`, `shred`, `chmod`, `chown`, `Remove-Item`, `del`, `rd`, `rmdir`, `format`, `Format-Volume`).

- **Alertas**: Quando um comando perigoso é sugerido, um alerta visual de aviso é exibido imediatamente abaixo do cartão de comando.
- **Palavra de Segurança**: Se você tentar executar um comando perigoso, será solicitado que você digite a palavra de segurança `"CONFIRM"` (diferencia maiúsculas de minúsculas) para prosseguir. Se for digitado incorretamente, a execução é abortada.
- **Lista Configurável**: A lista de comandos perigosos é totalmente editável.
- **Desativar Verificações de Segurança**: Você pode ignorar esta validação inteiramente, se desejar.

#### Comandos de Configuração:

```bash
# Adicionar/alterar a lista de comandos perigosos (lista separada por vírgulas)
shelloma config set dangerous "rm,dd,mkfs,shred,chmod,chown,Remove-Item,del,rd,rmdir,format,Format-Volume"

# Desativar completamente a verificação de comandos perigosos
shelloma config set disable_dangerous_check true

# Ativar a verificação novamente (padrão)
shelloma config set disable_dangerous_check false
```

### Flags Disponíveis

```text
Opções:
  -l, --lang string    Idioma: en, pt, es (padrão: en)
  -m, --model string   Modelo Ollama a utilizar (ex: qwen2.5-coder:1.5b)
  -url string          URL da API do Ollama (padrão: http://localhost:11434)
  -y, --yes            Executar o comando gerado automaticamente sem confirmação
  -v, --version        Exibir versão do Shelloma
```

---

## 🧪 Desenvolvimento & Qualidade de Código

O Shelloma adota práticas estritas de **Clean Code**, **Princípio de Responsabilidade Única (SRP)** e arquitetura organizada por funcionalidade (**Package by Feature**):

- **`pkg/cli`**: Orquestração da linha de comando, leitura de flags e fluxo interativo.
- **`pkg/ui`**: Interface visual de terminal, renderização de cards, estilos ANSI e área de transferência.
- **`pkg/ollama`**: Cliente da API do Ollama, limpeza de prompts e diagnósticos de execução.
- **`pkg/sysinfo`**: Detecção de sistema operacional, distribuição Linux, shell e usuário.
- **`pkg/config`**: Gestão e persistência das configurações do usuário.
- **`pkg/i18n`**: Suporte nativo e embutido para internacionalização (inglês, português e espanhol).

### Análise Estática & Automação

```bash
# Executar suíte de testes unitários
make test

# Executar análise estática de código (golangci-lint / staticcheck / go vet)
make lint

# Compilar binário local (executa lint e test antes da compilação)
make build
```

O projeto inclui integração com **`golangci-lint`** e um **Git Pre-Commit Hook** (`.git/hooks/pre-commit`) que executa automaticamente a verificação estática de código e os testes unitários a cada commit e build.

---

## 📄 Licença

Este projeto está licenciado sob a licença **MIT** - consulte o arquivo [LICENSE](LICENSE) para obter mais detalhes.
