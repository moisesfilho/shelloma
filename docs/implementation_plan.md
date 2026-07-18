# Plano de Implementação: Shelloma (Tradutor de Linguagem Natural para Shell com Ollama)

O **Shelloma** é uma aplicação CLI nativa (compilada sem dependências runtime como Python ou Node.js) desenvolvida especificamente para distribuições Linux baseadas em Debian/Ubuntu. Ela se conecta ao Ollama local para traduzir instruções em linguagem natural para comandos do terminal Bash/Zsh com execução interativa e segura.

## User Review Required

> [!IMPORTANT]
> **Linguagem Recomendada: Go (Golang)**
> Recomendamos **Go (Golang)** para este projeto pelos seguintes motivos:
> 1. Compila para um binário nativo estático/dinâmico sem dependências no sistema do usuário.
> 2. Possui cliente HTTP e suporte a JSON nativos e ultrarrápidos na stdlib.
> 3. É a linguagem padrão da indústria para CLI Linux (`docker`, `gh`, `k9s`).
> 4. Facilidade para gerar o pacote Debian `.deb` distribuível via `dpkg-deb`.
> 
> *Nota*: O repositório atual `/home/moises/Projetos/shelloma` está vazio e pronto para receber o código.

> [!TIP]
> **Modelos Ollama Detectados**:
> Identificamos que você já possui os seguintes modelos instalados no seu Ollama:
> - `qwen2.5-coder:1.5b` *(Recomendado para início por ser rápido e especializado em código/shell)*
> - `deepseek-coder:latest`
> - `llama3.2:latest`
> - `deepseek-r1:8b`

## Open Questions

> [!NOTE]
> 1. **Instalação do Compilador Go**: Para compilar o projeto no seu ambiente, podemos instalar o pacote `golang-go` via `apt` ou utilizar o Go oficial. Deseja que instalemos o pacote Go via `apt` durante a execução?
> 2. **Nome do Comando**: O nome do executável CLI será `shelloma` (exemplo de uso: `shelloma "listar arquivos grandes"`). Confirma esse nome ou prefere um alias (como `shll` ou `oma`)?

## Proposed Changes

Criaremos a estrutura completa do projeto nativo em Go em `/home/moises/Projetos/shelloma`.

---

### [Shelloma CLI App Core]

#### [NEW] [go.mod](file:///home/moises/Projetos/shelloma/go.mod)
- Módulo Go para o projeto `shelloma`.

#### [NEW] [main.go](file:///home/moises/Projetos/shelloma/main.go)
- Ponto de entrada do executável.
- Tratamento de sinalizadores (flags) como `--model`, `--url`, `--yes` (execução automática sem prompt), `config` e ajuda.
- Leitura dos argumentos de linguagem natural.

#### [NEW] [pkg/config/config.go](file:///home/moises/Projetos/shelloma/pkg/config/config.go)
- Carregamento e salvamento de configurações em `~/.config/shelloma/config.json`.
- Configurações suportadas:
  - `ollama_url`: URL do serviço Ollama (Padrão: `http://localhost:11434`)
  - `model`: Modelo utilizado (Padrão: auto-detectado entre `qwen2.5-coder`, `deepseek-coder`, `llama3.2` ou configurado pelo usuário)
  - `temperature`: Ajuste de aleatoriedade (Padrão: `0.1` para respostas precisas)
  - `auto_execute`: Booleano para execução direta sem confirmação.

#### [NEW] [pkg/sysinfo/sysinfo.go](file:///home/moises/Projetos/shelloma/pkg/sysinfo/sysinfo.go)
- Coleta contexto do sistema do usuário para enriquecer o prompt da IA:
  - Distribuição Linux (lendo `/etc/os-release`, ex: Ubuntu 24.04 LTS)
  - Shell em uso (`$SHELL`, ex: `/bin/bash` ou `/bin/zsh`)
  - Diretório atual (`$PWD`)
  - Usuário atual (`$USER` e se possui `sudo` / root)
  - Arquitetura (`amd64`/`arm64`)

#### [NEW] [pkg/ollama/client.go](file:///home/moises/Projetos/shelloma/pkg/ollama/client.go)
- Cliente HTTP nativo para se comunicar com a REST API do Ollama (`/api/generate` e `/api/tags`).
- Validação se o serviço Ollama está rodando.
- Construção do System Prompt otimizado para gerar **apenas o comando shell puro** (sem blocos de código markdown ou explicações desnecessárias no retorno principal).
- Suporte a requisições de explicação (`Explain`) do comando gerado.

#### [NEW] [pkg/ui/ui.go](file:///home/moises/Projetos/shelloma/pkg/ui/ui.go)
- Interface de terminal interativa com cores ANSI:
  - Caixa estilizada destacando o comando sugerido.
  - Opções interativas do teclado:
    - **`[Enter]` / `[y]` Execute**: Executa o comando diretamente no terminal via subshell e exibe a saída.
    - **`[e]` Explain**: Solicita ao Ollama uma explicação detalhada linha a linha do que o comando faz.
    - **`[m]` Edit**: Permite ao usuário editar o texto do comando antes de rodar.
    - **`[c]` Copy**: Copia o comando para a área de transferência do sistema (com suporte a OSC 52 e `xclip`/`wl-copy`).
    - **`[q]` / `[n]` Quit**: Aborta a operação sem executar nada.

---

### [Packaging & Automation]

#### [NEW] [Makefile](file:///home/moises/Projetos/shelloma/Makefile)
- Targets para:
  - `make build`: Compila o binário estático `shelloma`.
  - `make install`: Compila e instala em `/usr/local/bin/shelloma`.
  - `make deb`: Gera o pacote Debian `.deb` pronto para distribuição em Ubuntu/Debian.
  - `make clean`: Limpa artefatos de build.

#### [NEW] [scripts/build-deb.sh](file:///home/moises/Projetos/shelloma/scripts/build-deb.sh)
- Script automatizado para criar a estrutura do pacote Debian:
  - `DEBIAN/control`: Metadados do pacote (nome, versão, arquitetura, dependência recomendada de `ollama`).
  - Instalação do binário `/usr/bin/shelloma`.
  - Script de autocompletar para Bash/Zsh.

---

## Verification Plan

### Automated Tests
- Compilação limpa via Go (`go build -o shelloma .`).
- Testes unitários para parsing de respostas e sysinfo (`go test ./...`).

### Manual Verification
1. **Verificação de conexão Ollama**: Testar se o `shelloma` detecta os modelos locais (`qwen2.5-coder:1.5b`, `deepseek-coder`, etc.).
2. **Tradução e Execução**:
   - `shelloma "listar os 5 maiores diretórios no caminho atual"`
   - `shelloma "qual é meu IP local"`
   - `shelloma "procurar arquivos .log modificados nas últimas 24 horas"`
3. **Menu Interativo**: Testar as ações Execute, Explain, Edit e Cancel.
4. **Gerador de Pacote `.deb`**: Criar e instalar o pacote `.deb` usando `sudo dpkg -i` para verificar o processo de instalação em distribuições Debian/Ubuntu.
