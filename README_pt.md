# 🐚 Shelloma

> **Tradutor de Linguagem Natural para Comandos Terminal (Linux, macOS e Windows) alimentado por Ollama local.**

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![CI](https://github.com/moisesfilho/shelloma/actions/workflows/ci.yml/badge.svg)](https://github.com/moisesfilho/shelloma/actions/workflows/ci.yml)
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
- 🌐 **Internacionalização Nativa (i18n)**: Suporte completo para **Português Brasileiro (`pt`)**, **Inglês (`en`)** e **Espanhol (`es`)** através de arquivos JSON embarcados (com cabeçalho e mensagens traduzidos).
- 🤖 **Seleção Automática de Modelo**: Detecta os modelos instalados no Ollama e seleciona automaticamente o melhor modelo focado em código/shell disponível.
- 💡 **Detecção e Inicialização do Ollama**: Se o serviço Ollama estiver parado, o Shelloma avisa e oferece um comando rápido interativo para iniciá-lo de acordo com seu SO.
- 🖥️ **Interface Adaptativa com Rodapé Fixo**: Apenas o rodapé com a legenda de atalhos permanece fixado na base do terminal, sem atrapalhar a rolagem natural da zona de digitação e do cabeçalho.
- 🔄 **Refinamento Interativo de Comandos**: Permite ajustar o comando sugerido escolhendo a opção `[r: Refinar]`, abrindo uma caixa interativa para complementar a instrução original.
- 🧭 **Histórico Interativo com Persistência**: Todos os prompts interativos usam um leitor raw que habilita navegação no histórico com setas `Up`/`Down` e posicionamento por setas `Left`/`Right`. O histórico é persistido entre sessões (até 1000 entradas).
- 🖥️ **Modo Terminal Desktop**: Execute com a flag `--desktop` para manter o terminal aberto após a execução ("Pressione Enter para sair..."). Ideal para lançadores e integração com arquivos `.desktop`.
- ♾️ **Loop Contínuo de Comandos**: Após uma execução bem-sucedida, o Shelloma pergunta se deseja executar outra instrução, permitindo um fluxo contínuo sem reiniciar.
- ⌨️ **Atalhos de Teclado Avançados**: O menu interativo suporta múltiplos aliases por ação (ex: `y`/`sim`/`yes` para executar, `e`/`ex`/`explain` para explicar), além de `Esc` para sair de qualquer prompt.
- 📦 **Instalação Multiplataforma**: Executáveis compilados (`.exe`, binários nativos) e pacotes Linux (`.deb`, `.rpm`, Snap, AppImage, Flatpak). Inclui integração desktop Linux (gerando o atalho `.desktop` e registrando o ícone oficial do aplicativo em formato vetorial SVG).
- ⛓️ **Execução Multi Etapa**: Executa sequências de comandos multi etapa passo a passo com confirmações e gera logs para cada comando.
- 📝 **Completion Bash**: Scripts de completude para shell automáticos em pacotes Debian/Ubuntu.
- 🧠 **Aprendizado de Comandos**: Ensine um comando CLI ao Shelloma localmente com `shelloma learn <cmd>` (alias `aprender`). O help aprendido é salvo na sua máquina e injetado automaticamente no prompt sempre que o comando for mencionado.

---

## ⚙️ Dependências

1. **Sistema Operacional**: Linux, macOS ou Windows.
2. **Ollama**: Serviço do Ollama instalado e ativo localmente.
3. **Go 1.23+** *(opcional)*: Apenas caso deseje compilar a aplicação a partir do código-fonte.

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

Todos os binários e pacotes pré-compilados do Shelloma estão disponíveis para download na página oficial de **[Releases no GitHub](https://github.com/moisesfilho/shelloma/releases)** (`.deb`, `.rpm`, `.tar.gz`, `.zip`, `.AppImage`, `.flatpak`, Snap, `.exe`).

### Opção 1: Download de Binários e Pacotes Pré-compilados (Recomendado)

Acesse a página de **[GitHub Releases](https://github.com/moisesfilho/shelloma/releases)** e escolha a opção ideal para o seu sistema operacional:

#### 🐧 Linux
- **Debian / Ubuntu / Linux Mint / Pop!_OS (`.deb`)**:
  ```bash
  sudo dpkg -i shelloma_1.3.0_amd64.deb
  ```
  *(Durante a instalação do `.deb`, será exibido um assistente no terminal para selecionar o idioma padrão).*

- **Fedora / RedHat / CentOS / RHEL (`.rpm`)**:
  ```bash
  sudo rpm -i shelloma_1.3.0_amd64.rpm
  ```

- **Flatpak (Universal)**:
  ```bash
  flatpak install Shelloma-x86_64.flatpak
  ```

- **Snap (Universal)**:
  ```bash
  snap install shelloma
  ```
  *(Ou baixe o `Shelloma-x86_64.snap` das Releases e instale com `snap install --dangerous Shelloma-x86_64.snap`)*

- **AppImage (Universal - Executável portável sem necessidade de instalação)**:
  ```bash
  chmod +x Shelloma-x86_64.AppImage
  ./Shelloma-x86_64.AppImage "listar arquivos do diretorio"
  ```

- **Arquivo Comprimido Tarball (`.tar.gz`)**:
  Baixe o `.tar.gz` (`amd64` ou `arm64`), extraia o binário e mova para seu `PATH` (ex: `~/.local/bin/`).

#### 🍏 macOS
- **Binário Nativo (Intel & Apple Silicon M1/M2/M3)**:
  1. Baixe `shelloma_1.3.0_darwin_arm64.tar.gz` (Apple Silicon) ou `shelloma_1.3.0_darwin_amd64.tar.gz` (Intel) na página de Releases.
  2. Extraia o pacote e torne-o executável:
     ```bash
     tar -xzf shelloma_1.3.0_darwin_arm64.tar.gz
     chmod +x shelloma
     sudo mv shelloma /usr/local/bin/
     ```

#### 🪟 Windows
- **Executável Nativamente Compilado (`.exe` / `.zip`)**:
  1. Baixe `shelloma_1.3.0_windows_amd64.zip` (64-bit) ou `shelloma_1.3.0_windows_arm64.zip` (ARM64) na página de Releases.
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

Ao passar uma instrução diretamente, o Shelloma executa em **modo execução única**: sugere, executa e encerra — sem o prompt de continuação do loop. Execute `shelloma` sem argumentos para usar o app interativo completo.

### Menu de Opções da CLI

Ao gerar o comando, o Shelloma exibirá o card e aguardará sua ação:

```text
Opções: [Enter/y: Executar] [e: Explicar] [m: Modificar] [c: Copiar] [p: Novo Prompt] [r: Refinar] [a: Ajustar Prompt] [q/n: Sair]:
```

- **Pressionar `Enter` ou `y`**: Executa o comando diretamente no terminal.
- **Digite `e`**: Solicita uma explicação detalhada e linha por linha do comando ao Ollama.
- **Digite `m`**: Abre um prompt para você reescrever/modificar o comando manualmente antes de executar.
- **Digite `c`**: Copia o comando diretamente para a área de transferência do sistema (clipboard).
- **Digite `p`**: Solicita um novo prompt de comandos do início, sem precisar sair do Shelloma.
- **Digite `r`**: Refina o comando gerado, fornecendo feedback ou ajustes adicionais para que o LLM crie uma nova sugestão revisada.
- **Digite `a`**: Permite reajustar/editar o prompt textual original que deu origem à sugestão de comando.
- **Digite `q` ou `n`**: Cancela a operação e encerra a aplicação com segurança.
- **Pressione `Esc`**: Sai imediatamente de qualquer prompt ou menu interativo.

### Comandos de Configuração e Utilitários

```bash
# Exibir as configurações atuais do Shelloma
shelloma config

# Exibir a documentação compactada das configurações (schema)
shelloma config docs

# Alterar o idioma padrão salvo (pt, en, es)
shelloma config set lang pt

# Definir um modelo específico do Ollama
shelloma config set model qwen2.5-coder:1.5b

# Definir a temperatura de amostragem do modelo (padrão: 0.1)
shelloma config set temperature 0.5

# Ativar/Desativar execução automática de comandos sem confirmação
shelloma config set auto_execute true

# Listar os modelos do Ollama instalados na sua máquina
shelloma models

# Abrir os logs de execução da aplicação
shelloma logs
```

### 🤖 Configuração por Linguagem Natural

O Shelloma permite que você solicite atualizações de configuração por meio de linguagem natural! A IA local entende o pedido e gera o comando CLI do Shelloma correspondente para ser executado:

- `shelloma "mude o modelo do shelloma para llama3"` ➔ Sugere: `shelloma config set model llama3`
- `shelloma "mude o idioma para espanhol"` ➔ Sugere: `shelloma config set lang es`
- `shelloma "defina a temperatura como 0.8"` ➔ Sugere: `shelloma config set temperature 0.8`
- `shelloma "ative a execução automática"` ➔ Sugere: `shelloma config set auto_execute true`

### 📝 Regras Customizadas do Provedor (Rules)

Você pode definir regras de formatação e preferências de execução personalizadas que serão injetadas diretamente no prompt de sistema do Ollama para orientar a geração de comandos (como definir pastas padrão, escolher editores de texto CLI ou programas específicos para abrir determinadas extensões de arquivos).

Comandos:
```bash
# Adicionar uma nova regra customizada
shelloma rules add "Sempre abrir imagens com xdg-open"

# Listar todas as regras salvas
shelloma rules list

# Editar uma regra salva pelo índice
shelloma rules edit 1 "Sempre abrir imagens com feh"

# Excluir uma regra salva pelo índice
shelloma rules delete 1
```

### 🧠 Aprendizado de Comandos Customizados

O Shelloma pode aprender qualquer comando CLI disponível na sua máquina para que futuras solicitações que o envolvam sejam respondidas com as opções reais do comando:

```bash
# Aprender as opções de um comando localmente (alias: aprender)
shelloma learn git
shelloma aprender git
```

Os comandos aprendidos ficam salvos localmente em `~/.config/shelloma/learned/<comando>.json`. Quando você menciona um comando aprendido na sua solicitação, a referência de help é injetada automaticamente no prompt para que o modelo gere comandos precisos.

### 📋 Histórico de Execução (Logs)

Toda sugestão de comando e resultado de execução é registrado automaticamente em um arquivo de log estruturado (utilizando o diretório de cache padrão de cada sistema operacional: `~/.cache/shelloma/shelloma.log` no Linux/macOS, ou `%LocalAppData%\shelloma\shelloma.log` no Windows).

Execute `shelloma logs` para abrir os logs. Você pode escolher:
1. Exibir os logs em um formato limpo e estruturado diretamente no seu terminal.
2. Abrir o arquivo de log no editor de texto padrão do seu sistema.

### ⛓️ Execução de Comando Multi Etapa

Quando uma instrução do usuário exige várias etapas de execução, o Shelloma instrui o modelo local a estruturar os comandos em linhas separadas (um comando por linha).

Recursos:
- **Confirmação por Etapa**: O Shelloma exibe a sequência planejada de comandos e solicita a confirmação do usuário antes de executar cada etapa individualmente.
- **Persistência de Navegação de Diretório**: Comandos `cd` são interceptados e aplicados ao diretório do processo, garantindo que as etapas subsequentes rodem no contexto correto.
- **Análise de Erros e Logs**: Cada etapa individual é validada contra falhas e gera sua própria entrada no histórico de logs.

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
      --desktop        Modo terminal desktop (pausa com "Pressione Enter para sair...")
  -v, --version        Exibir versão do Shelloma
```

> **Modo Desktop**: Ao iniciar por um arquivo `.desktop` ou lançador gráfico, use `--desktop` para manter o terminal aberto após a execução — útil para ver a saída antes da janela fechar.

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

# Gerar ícones da aplicação (requer Python + Pillow)
make icons

# Compilar binário local (executa lint, test e icons antes da compilação)
make build

# Compilar para todas as plataformas (Linux, macOS, Windows — amd64 + arm64)
make build-all

# Gerar pacotes .deb, .tar.gz, .AppImage ou .flatpak
make deb
make tar
make appimage
make flatpak

# Instalar binário + .desktop + ícones em ~/.local/bin/
make install-user
```

O projeto inclui integração com **`golangci-lint`**, um **Git Pre-Commit Hook** (`.git/hooks/pre-commit`) e um **pipeline CI/CD** via GitHub Actions:
- **CI Quality Gate** (`.github/workflows/ci.yml`): Executa lint + test em push/PR para `main` e `develop`.
- **Release Pipeline** (`.github/workflows/release.yml`): Ao enviar uma tag, executa lint + test e publica todos os pacotes (`.deb`, `.rpm`, `.tar.gz`, `.zip`, `.AppImage`, `.flatpak`, Snap) no GitHub Releases.

---

## 📄 Licença

Este projeto está licenciado sob a licença **MIT** - consulte o arquivo [LICENSE](LICENSE) para obter mais detalhes.
